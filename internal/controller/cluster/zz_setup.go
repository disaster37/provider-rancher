// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	appv2 "github.com/disaster37/provider-rancher/internal/controller/app/appv2"
	catalogv2 "github.com/disaster37/provider-rancher/internal/controller/app/catalogv2"
	authconfigactivedirectory "github.com/disaster37/provider-rancher/internal/controller/auth/authconfigactivedirectory"
	authconfigadfs "github.com/disaster37/provider-rancher/internal/controller/auth/authconfigadfs"
	authconfigazuread "github.com/disaster37/provider-rancher/internal/controller/auth/authconfigazuread"
	authconfigfreeipa "github.com/disaster37/provider-rancher/internal/controller/auth/authconfigfreeipa"
	authconfiggithub "github.com/disaster37/provider-rancher/internal/controller/auth/authconfiggithub"
	authconfigkeycloak "github.com/disaster37/provider-rancher/internal/controller/auth/authconfigkeycloak"
	authconfigokta "github.com/disaster37/provider-rancher/internal/controller/auth/authconfigokta"
	authconfigopenldap "github.com/disaster37/provider-rancher/internal/controller/auth/authconfigopenldap"
	authconfigping "github.com/disaster37/provider-rancher/internal/controller/auth/authconfigping"
	certificate "github.com/disaster37/provider-rancher/internal/controller/k8s/certificate"
	configmapv2 "github.com/disaster37/provider-rancher/internal/controller/k8s/configmapv2"
	project "github.com/disaster37/provider-rancher/internal/controller/k8s/project"
	ranchernamespace "github.com/disaster37/provider-rancher/internal/controller/k8s/ranchernamespace"
	registry "github.com/disaster37/provider-rancher/internal/controller/k8s/registry"
	secretv2 "github.com/disaster37/provider-rancher/internal/controller/k8s/secretv2"
	storageclassv2 "github.com/disaster37/provider-rancher/internal/controller/k8s/storageclassv2"
	providerconfig "github.com/disaster37/provider-rancher/internal/controller/providerconfig"
	bootstrap "github.com/disaster37/provider-rancher/internal/controller/rancher/bootstrap"
	cloudcredential "github.com/disaster37/provider-rancher/internal/controller/rancher/cloudcredential"
	cluster "github.com/disaster37/provider-rancher/internal/controller/rancher/cluster"
	clusterdriver "github.com/disaster37/provider-rancher/internal/controller/rancher/clusterdriver"
	clusterroletemplatebinding "github.com/disaster37/provider-rancher/internal/controller/rancher/clusterroletemplatebinding"
	clustersync "github.com/disaster37/provider-rancher/internal/controller/rancher/clustersync"
	clustertemplate "github.com/disaster37/provider-rancher/internal/controller/rancher/clustertemplate"
	clusterv2 "github.com/disaster37/provider-rancher/internal/controller/rancher/clusterv2"
	customusertoken "github.com/disaster37/provider-rancher/internal/controller/rancher/customusertoken"
	etcdbackup "github.com/disaster37/provider-rancher/internal/controller/rancher/etcdbackup"
	feature "github.com/disaster37/provider-rancher/internal/controller/rancher/feature"
	globalrole "github.com/disaster37/provider-rancher/internal/controller/rancher/globalrole"
	globalrolebinding "github.com/disaster37/provider-rancher/internal/controller/rancher/globalrolebinding"
	machineconfigv2 "github.com/disaster37/provider-rancher/internal/controller/rancher/machineconfigv2"
	nodedriver "github.com/disaster37/provider-rancher/internal/controller/rancher/nodedriver"
	nodepool "github.com/disaster37/provider-rancher/internal/controller/rancher/nodepool"
	nodetemplate "github.com/disaster37/provider-rancher/internal/controller/rancher/nodetemplate"
	roletemplate "github.com/disaster37/provider-rancher/internal/controller/rancher/roletemplate"
	setting "github.com/disaster37/provider-rancher/internal/controller/rancher/setting"
	token "github.com/disaster37/provider-rancher/internal/controller/rancher/token"
	user "github.com/disaster37/provider-rancher/internal/controller/rancher/user"
	projectroletemplatebinding "github.com/disaster37/provider-rancher/internal/controller/rbac/projectroletemplatebinding"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		appv2.Setup,
		catalogv2.Setup,
		authconfigactivedirectory.Setup,
		authconfigadfs.Setup,
		authconfigazuread.Setup,
		authconfigfreeipa.Setup,
		authconfiggithub.Setup,
		authconfigkeycloak.Setup,
		authconfigokta.Setup,
		authconfigopenldap.Setup,
		authconfigping.Setup,
		certificate.Setup,
		configmapv2.Setup,
		project.Setup,
		ranchernamespace.Setup,
		registry.Setup,
		secretv2.Setup,
		storageclassv2.Setup,
		providerconfig.Setup,
		bootstrap.Setup,
		cloudcredential.Setup,
		cluster.Setup,
		clusterdriver.Setup,
		clusterroletemplatebinding.Setup,
		clustersync.Setup,
		clustertemplate.Setup,
		clusterv2.Setup,
		customusertoken.Setup,
		etcdbackup.Setup,
		feature.Setup,
		globalrole.Setup,
		globalrolebinding.Setup,
		machineconfigv2.Setup,
		nodedriver.Setup,
		nodepool.Setup,
		nodetemplate.Setup,
		roletemplate.Setup,
		setting.Setup,
		token.Setup,
		user.Setup,
		projectroletemplatebinding.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
