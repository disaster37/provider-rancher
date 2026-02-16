/*
Copyright 2021 Upbound Inc.
*/

package clients

import (
	"context"
	"encoding/json"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/upjet/v2/pkg/terraform"

	"github.com/disaster37/provider-rancher/apis/v1beta1"
)

const (
	// error messages
	errNoProviderConfig     = "no providerConfigRef provided"
	errGetProviderConfig    = "cannot get referenced ProviderConfig"
	errTrackUsage           = "cannot track ProviderConfig usage"
	errExtractCredentials   = "cannot extract credentials"
	errUnmarshalCredentials = "cannot unmarshal rancher credentials as JSON"
	keyApiURL               = "api_url"
	keyAccessKey            = "access_key"
	keySecretKey            = "secret_key"
	keyBootstrap            = "bootstrap"
	keyAlias                = "alias"
	keyTokenKey             = "token_key"
	keyInsecure             = "insecure"
	keyCACerts              = "ca_certs"
	keyRetries              = "retries"
	keyTimeout              = "timeout"
)

// TerraformSetupBuilder builds Terraform a terraform.SetupFn function which
// returns Terraform provider setup configuration
func TerraformSetupBuilder(version, providerSource, providerVersion string) terraform.SetupFn {
	return func(ctx context.Context, client client.Client, mg resource.Managed) (terraform.Setup, error) {
		ps := terraform.Setup{
			Version: version,
			Requirement: terraform.ProviderRequirement{
				Source:  providerSource,
				Version: providerVersion,
			},
		}

		configRef := mg.GetProviderConfigReference()
		if configRef == nil {
			return ps, errors.New(errNoProviderConfig)
		}
		pc := &v1beta1.ProviderConfig{}
		if err := client.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc); err != nil {
			return ps, errors.Wrap(err, errGetProviderConfig)
		}

		t := resource.NewProviderConfigUsageTracker(client, &v1beta1.ProviderConfigUsage{})
		if err := t.Track(ctx, mg); err != nil {
			return ps, errors.Wrap(err, errTrackUsage)
		}

		data, err := resource.CommonCredentialExtractor(ctx, pc.Spec.Credentials.Source, client, pc.Spec.Credentials.CommonCredentialSelectors)
		if err != nil {
			return ps, errors.Wrap(err, errExtractCredentials)
		}
		creds := map[string]string{}
		if err := json.Unmarshal(data, &creds); err != nil {
			return ps, errors.Wrap(err, errUnmarshalCredentials)
		}

		// set provider configuration
		ps.Configuration = map[string]any{}
		if v, ok := creds[keyApiURL]; ok {
			ps.Configuration[keyApiURL] = v
		}
		if v, ok := creds[keyAccessKey]; ok {
			ps.Configuration[keyAccessKey] = v
		}
		if v, ok := creds[keySecretKey]; ok {
			ps.Configuration[keySecretKey] = v
		}
		if v, ok := creds[keyAlias]; ok {
			ps.Configuration[keyAlias] = v
		}
		if v, ok := creds[keyBootstrap]; ok {
			ps.Configuration[keyBootstrap] = v
		}
		if v, ok := creds[keyCACerts]; ok {
			ps.Configuration[keyCACerts] = v
		}
		if v, ok := creds[keyInsecure]; ok {
			ps.Configuration[keyInsecure] = v
		}
		if v, ok := creds[keyRetries]; ok {
			ps.Configuration[keyRetries] = v
		}
		if v, ok := creds[keyTimeout]; ok {
			ps.Configuration[keyTimeout] = v
		}
		if v, ok := creds[keyTokenKey]; ok {
			ps.Configuration[keyTokenKey] = v
		}

		return ps, nil
	}
}
