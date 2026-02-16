/*
Copyright 2021 Upbound Inc.
*/

package main

import (
	"os"
	"path/filepath"
	"time"

	xpcontroller "github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/feature"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/statemetrics"
	tjcontroller "github.com/crossplane/upjet/v2/pkg/controller"
	"github.com/crossplane/upjet/v2/pkg/terraform"
	"gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/disaster37/provider-rancher/apis"
	"github.com/disaster37/provider-rancher/config"
	"github.com/disaster37/provider-rancher/internal/clients"
	"github.com/disaster37/provider-rancher/internal/controller"
	"github.com/disaster37/provider-rancher/internal/features"
)

func main() {
	var (
		app                     = kingpin.New(filepath.Base(os.Args[0]), "Terraform based Crossplane provider for Rancher").DefaultEnvars()
		debug                   = app.Flag("debug", "Run with debug logging.").Short('d').Bool()
		syncPeriod              = app.Flag("sync", "Controller manager sync period such as 300ms, 1.5h, or 2h45m").Short('s').Default("1h").Duration()
		pollInterval            = app.Flag("poll", "Poll interval controls how often an individual resource should be checked for drift.").Default("10m").Duration()
		pollStateMetricInterval = app.Flag("poll-state-metric", "State metric recording interval").Default("5s").Duration()
		leaderElection          = app.Flag("leader-election", "Use leader election for the controller manager.").Short('l').Default("false").OverrideDefaultFromEnvar("LEADER_ELECTION").Bool()
		maxReconcileRate        = app.Flag("max-reconcile-rate", "The global maximum rate per second at which resources may be checked for drift from the desired state.").Default("10").Int()

		terraformVersion = app.Flag("terraform-version", "Terraform version.").Required().Envar("TERRAFORM_VERSION").String()
		providerSource   = app.Flag("terraform-provider-source", "Terraform provider source.").Required().Envar("TERRAFORM_PROVIDER_SOURCE").String()
		providerVersion  = app.Flag("terraform-provider-version", "Terraform provider version.").Required().Envar("TERRAFORM_PROVIDER_VERSION").String()

		enableManagementPolicies = app.Flag("enable-management-policies", "Enable support for Management Policies.").Default("true").Envar("ENABLE_MANAGEMENT_POLICIES").Bool()
	)

	kingpin.MustParse(app.Parse(os.Args[1:]))

	zl := zap.New(zap.UseDevMode(*debug))
	log := logging.NewLogrLogger(zl.WithName("provider-rancher"))
	if *debug {
		// The controller-runtime runs with a no-op logger by default. It is
		// *very* verbose even at info level, so we only provide it a real
		// logger when we're running in debug mode.
		ctrl.SetLogger(zl)
	}

	log.Debug("Starting", "sync-period", syncPeriod.String(), "poll-interval", pollInterval.String(), "max-reconcile-rate", *maxReconcileRate)

	cfg, err := ctrl.GetConfig()
	kingpin.FatalIfError(err, "Cannot get API server rest config")

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		LeaderElection:   *leaderElection,
		LeaderElectionID: "crossplane-leader-election-provider-rancher",
		Cache: cache.Options{
			SyncPeriod: syncPeriod,
		},
		LeaderElectionResourceLock: resourcelock.LeasesResourceLock,
		LeaseDuration:              func() *time.Duration { d := 60 * time.Second; return &d }(),
		RenewDeadline:              func() *time.Duration { d := 50 * time.Second; return &d }(),
	})
	kingpin.FatalIfError(err, "Cannot create controller manager")
	kingpin.FatalIfError(apis.AddToScheme(mgr.GetScheme()), "Cannot add Rancher APIs to scheme")

	metricRecorder := managed.NewMRMetricRecorder()
	stateMetrics := statemetrics.NewMRStateMetrics()

	metrics.Registry.MustRegister(metricRecorder)
	metrics.Registry.MustRegister(stateMetrics)

	o := tjcontroller.Options{
		Options: xpcontroller.Options{
			Logger:                  log,
			GlobalRateLimiter:       ratelimiter.NewGlobal(*maxReconcileRate),
			PollInterval:            *pollInterval,
			MaxConcurrentReconciles: *maxReconcileRate,
			Features:                &feature.Flags{},
			MetricOptions: &xpcontroller.MetricOptions{
				PollStateMetricInterval: *pollStateMetricInterval,
				MRMetrics:               metricRecorder,
				MRStateMetrics:          stateMetrics,
			},
		},
		Provider: config.GetProvider(),
		// use the following WorkspaceStoreOption to enable the shared gRPC mode
		// terraform.WithProviderRunner(terraform.NewSharedProvider(log, os.Getenv("TERRAFORM_NATIVE_PROVIDER_PATH"), terraform.WithNativeProviderArgs("-debuggable")))
		WorkspaceStore: terraform.NewWorkspaceStore(log),
		SetupFn:        clients.TerraformSetupBuilder(*terraformVersion, *providerSource, *providerVersion),
	}

	if *enableManagementPolicies {
		o.Features.Enable(features.EnableBetaManagementPolicies)
		log.Info("Beta feature enabled", "flag", features.EnableBetaManagementPolicies)
	}

	kingpin.FatalIfError(controller.Setup(mgr, o), "Cannot setup Rancher controllers")
	kingpin.FatalIfError(mgr.Start(ctrl.SetupSignalHandler()), "Cannot start controller manager")
}
