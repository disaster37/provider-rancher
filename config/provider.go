/*
Copyright 2021 Upbound Inc.
*/

package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"
	"strings"

	"github.com/crossplane/upjet/pkg/config"
	ujconfig "github.com/crossplane/upjet/pkg/config"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	resourcePrefix = "rancher"
	modulePath     = "github.com/disaster37/provider-rancher"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns provider configuration
func GetProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("rancher.contrib.crossplane.io"),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
			KindOverrides(),
		))

	for _, configure := range []func(provider *ujconfig.Provider){
		// add custom config functions
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_app", func(r *config.Resource) {
				r.ShortGroup = "app"
				r.MarkAsRequired(
					"catalog_name",
					"name",
					"target_namespace",
					"template_name",
				)

				r.References["catalog_name"] = config.Reference{
					TerraformName: "rancher2_catalog",
				}
				r.References["project_id"] = config.Reference{
					TerraformName: "rancher2_project",
				}
				r.References["target_namespace"] = config.Reference{
					TerraformName: "rancher2_namespace",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_app_v2", func(r *config.Resource) {
				r.ShortGroup = "app"
				r.Kind = "AppV2"
				r.MarkAsRequired(
					"name",
					"namespace",
					"repo_name",
					"chart_name",
				)

				r.References["cluster_id"] = config.Reference{
					TerraformName: "rancher2_cluster_v2",
				}
				r.References["project_id"] = config.Reference{
					TerraformName: "rancher2_project",
				}
				r.References["namespace"] = config.Reference{
					TerraformName: "rancher2_namespace",
				}
				r.References["repo_name"] = config.Reference{
					TerraformName: "rancher2_catalog_v2",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_auth_config_activedirectory", func(r *config.Resource) {
				r.ShortGroup = "auth"

				r.MarkAsRequired(
					"servers",
					"service_account_username",
					"service_account_password",
					"user_search_base",
					"test_username",
					"test_password",
				)
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_auth_config_adfs", func(r *config.Resource) {
				r.ShortGroup = "auth"

				r.MarkAsRequired(
					"display_name_field",
					"groups_field",
					"idp_metadata_content",
					"rancher_api_host",
					"sp_cert",
					"sp_key",
					"uid_field",
					"user_name_field",
				)
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_auth_config_azuread", func(r *config.Resource) {
				r.ShortGroup = "auth"

				r.MarkAsRequired(
					"application_id",
					"application_secret",
					"auth_endpoint",
					"graph_endpoint",
					"rancher_url",
					"tenant_id",
					"token_endpoint",
				)
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_auth_config_freeipa", func(r *config.Resource) {
				r.ShortGroup = "auth"

				r.MarkAsRequired(
					"servers",
					"service_account_distinguished_name",
					"service_account_password",
					"user_search_base",
					"test_username",
					"test_password",
				)
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_auth_config_github", func(r *config.Resource) {
				r.ShortGroup = "auth"

				r.MarkAsRequired(
					"client_id",
					"client_secret",
				)
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_auth_config_keycloak", func(r *config.Resource) {
				r.ShortGroup = "auth"

				r.MarkAsRequired(
					"display_name_field",
					"groups_field",
					"idp_metadata_content",
					"rancher_api_host",
					"sp_cert",
					"sp_key",
					"uid_field",
					"user_name_field",
				)
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_auth_config_okta", func(r *config.Resource) {
				r.ShortGroup = "auth"

				r.MarkAsRequired(
					"display_name_field",
					"groups_field",
					"idp_metadata_content",
					"rancher_api_host",
					"sp_cert",
					"sp_key",
					"uid_field",
					"user_name_field",
				)
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_auth_config_openldap", func(r *config.Resource) {
				r.ShortGroup = "auth"

				r.MarkAsRequired(
					"servers",
					"service_account_distinguished_name",
					"service_account_password",
					"user_search_base",
					"test_username",
					"test_password",
				)
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_auth_config_ping", func(r *config.Resource) {
				r.ShortGroup = "auth"

				r.MarkAsRequired(
					"display_name_field",
					"groups_field",
					"idp_metadata_content",
					"rancher_api_host",
					"sp_cert",
					"sp_key",
					"uid_field",
					"user_name_field",
				)
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_bootstrap", func(r *config.Resource) {
				r.ShortGroup = "rancher"
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_catalog", func(r *config.Resource) {
				r.ShortGroup = "app"

				r.MarkAsRequired(
					"name",
					"url",
				)

				r.References["cluster_id"] = config.Reference{
					TerraformName: "rancher2_cluster",
				}
				r.References["project_id"] = config.Reference{
					TerraformName: "rancher2_project",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_catalog_v2", func(r *config.Resource) {
				r.ShortGroup = "app"
				r.Kind = "CatalogV2"

				r.MarkAsRequired(
					"name",
				)

				r.References["cluster_id"] = config.Reference{
					TerraformName: "rancher2_cluster_v2",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_certificate", func(r *config.Resource) {
				r.ShortGroup = "k8s"

				r.MarkAsRequired(
					"certs",
					"key",
				)

				r.References["namespace_id"] = config.Reference{
					TerraformName: "rancher2_namespace",
				}
				r.References["project_id"] = config.Reference{
					TerraformName: "rancher2_project",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_cloud_credential", func(r *config.Resource) {
				r.ShortGroup = "rancher"

				r.MarkAsRequired(
					"name",
				)
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_cluster", func(r *config.Resource) {
				r.ShortGroup = "rancher"

				r.MarkAsRequired(
					"name",
				)

				r.References["cluster_template_id"] = config.Reference{
					TerraformName: "rancher2_cluster_template",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_cluster_driver", func(r *config.Resource) {
				r.ShortGroup = "rancher"

				r.MarkAsRequired(
					"active",
					"builtin",
					"name",
					"url",
				)
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_cluster_role_template_binding", func(r *config.Resource) {
				r.ShortGroup = "rancher"

				r.MarkAsRequired(
					"name",
				)

				r.References["cluster_id"] = config.Reference{
					TerraformName: "rancher2_cluster_v2",
				}
				r.References["role_template_id"] = config.Reference{
					TerraformName: "rancher2_role_template",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_cluster_sync", func(r *config.Resource) {
				r.ShortGroup = "rancher"

				r.References["cluster_id"] = config.Reference{
					TerraformName: "rancher2_cluster_v2",
				}

				r.References["node_pool_ids"] = config.Reference{
					TerraformName: "rancher2_node_pool",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_cluster_template", func(r *config.Resource) {
				r.ShortGroup = "rancher"

				r.MarkAsRequired("name")
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_cluster_v2", func(r *config.Resource) {
				r.ShortGroup = "rancher"
				r.Kind = "ClusterV2"

				r.MarkAsRequired(
					"name",
					"kubernetes_version",
				)
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_config_map_v2", func(r *config.Resource) {
				r.ShortGroup = "k8s"
				r.Kind = "ConfigMapV2"

				r.MarkAsRequired(
					"name",
					"data",
				)

				r.References["cluster_id"] = config.Reference{
					TerraformName: "rancher2_cluster_v2",
				}
				r.References["namespace"] = config.Reference{
					TerraformName: "rancher2_namespace",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_custom_user_token", func(r *config.Resource) {
				r.ShortGroup = "rancher"

				r.MarkAsRequired(
					"username",
					"password",
				)

				r.References["username"] = config.Reference{
					TerraformName: "rancher2_user",
				}

			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_etcd_backup", func(r *config.Resource) {
				r.ShortGroup = "rancher"

				r.MarkAsRequired(
					"name",
				)

				r.References["cluster_id"] = config.Reference{
					TerraformName: "rancher2_cluster_v2",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_feature", func(r *config.Resource) {
				r.ShortGroup = "rancher"

				r.MarkAsRequired("name")
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_global_role", func(r *config.Resource) {
				r.ShortGroup = "rancher"

				r.MarkAsRequired("name")
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_global_role_binding", func(r *config.Resource) {
				r.ShortGroup = "rancher"

				r.References["global_role_id"] = config.Reference{
					TerraformName: "rancher2_global_role",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_machine_config_v2", func(r *config.Resource) {
				r.ShortGroup = "rancher"

				r.MarkAsRequired("generate_name")
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_multi_cluster_app", func(r *config.Resource) {
				r.ShortGroup = "app"

				r.MarkAsRequired(
					"catalog_name",
					"name",
					"targets",
					"template_name",
				)

				r.References["catalog_name"] = config.Reference{
					TerraformName: "rancher2_catalog_v2",
				}

				r.References["targets.project_id"] = config.Reference{
					TerraformName: "rancher2_project",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_namespace", func(r *config.Resource) {
				r.ShortGroup = "k8s"
				r.Kind = "RancherNamespace"

				r.MarkAsRequired(
					"name",
				)

				r.References["project_id"] = config.Reference{
					TerraformName: "rancher2_project",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_node_driver", func(r *config.Resource) {
				r.ShortGroup = "rancher"

				r.MarkAsRequired(
					"active",
					"builtin",
					"name",
					"url",
				)
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_node_pool", func(r *config.Resource) {
				r.ShortGroup = "rancher"

				r.MarkAsRequired(
					"name",
					"hostname_prefix",
				)

				r.References["cluster_id"] = config.Reference{
					TerraformName: "rancher2_cluster_v2",
				}

				r.References["node_template_id"] = config.Reference{
					TerraformName: "rancher2_node_template",
				}

				r.References["cloud_credential_id"] = config.Reference{
					TerraformName: "rancher2_cloud_credential",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_node_template", func(r *config.Resource) {
				r.ShortGroup = "rancher"

				r.MarkAsRequired("name")

				r.References["driver_id"] = config.Reference{
					TerraformName: "rancher2_node_driver",
				}

				r.References["cloud_credential_id"] = config.Reference{
					TerraformName: "rancher2_cloud_credential",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_project", func(r *config.Resource) {
				r.ShortGroup = "k8s"

				r.MarkAsRequired(
					"name",
				)

				r.References["cluster_id"] = config.Reference{
					TerraformName: "rancher2_cluster_v2",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_project_role_template_binding", func(r *config.Resource) {
				r.ShortGroup = "rbac"

				r.MarkAsRequired(
					"name",
				)

				r.References["project_id"] = config.Reference{
					TerraformName: "rancher2_cluster_v2",
				}

				r.References["role_template_id"] = config.Reference{
					TerraformName: "rancher2_role_template",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_registry", func(r *config.Resource) {
				r.ShortGroup = "k8s"

				r.MarkAsRequired(
					"name",
					"registries",
				)

				r.References["project_id"] = config.Reference{
					TerraformName: "rancher2_project",
				}

				r.References["namespace_id"] = config.Reference{
					TerraformName: "rancher2_namespace",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_role_template", func(r *config.Resource) {
				r.ShortGroup = "rancher"

				r.MarkAsRequired("name")
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_secret", func(r *config.Resource) {
				r.ShortGroup = "k8s"

				r.MarkAsRequired(
					"data",
				)

				r.References["project_id"] = config.Reference{
					TerraformName: "rancher2_project",
				}
				r.References["namespace_id"] = config.Reference{
					TerraformName: "rancher2_namespace",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_secret_v2", func(r *config.Resource) {
				r.ShortGroup = "k8s"
				r.Kind = "SecretV2"

				r.MarkAsRequired(
					"name",
					"data",
				)

				r.References["cluster_id"] = config.Reference{
					TerraformName: "rancher2_cluster_v2",
				}
				r.References["namespace"] = config.Reference{
					TerraformName: "rancher2_namespace",
				}

			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_setting", func(r *config.Resource) {
				r.ShortGroup = "rancher"

				r.MarkAsRequired(
					"name",
					"value",
				)
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_storage_class_v2", func(r *config.Resource) {
				r.ShortGroup = "k8s"
				r.Kind = "StorageClassV2"

				r.MarkAsRequired(
					"name",
					"k8s_provisioner",
				)

				r.References["cluster_id"] = config.Reference{
					TerraformName: "rancher2_cluster_v2",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_token", func(r *config.Resource) {
				r.ShortGroup = "rancher"

				r.References["cluster_id"] = config.Reference{
					TerraformName: "rancher2_cluster_v2",
				}
			})
		},
		func(p *ujconfig.Provider) {
			p.AddResourceConfigurator("rancher2_user", func(r *config.Resource) {
				r.ShortGroup = "rancher"

				r.MarkAsRequired(
					"username",
					"password",
				)
			})
		},
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}

// KindOverrides overrides the kind of the resources given in KindMap.
func KindOverrides() config.ResourceOption {
	return func(r *config.Resource) {
		subElement := strings.Split(strings.TrimPrefix(r.Name, "rancher2_"), "_")
		for i, elem := range subElement {
			subElement[i] = cases.Title(language.English).String(elem)
		}
		r.Kind = strings.Join(subElement, "")
	}
}
