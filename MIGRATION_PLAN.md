# Upjet v2 Migration Plan for provider-rancher

## Overview
This document outlines the step-by-step migration plan to upgrade provider-rancher from Upjet v1 to Upjet v2, enabling Crossplane v2 compatibility with namespaced Managed Resources.

## Current State Analysis

### Dependencies (from go.mod)
- `crossplane-runtime`: v1.16.0 → **needs upgrade to v2.0.0**
- `crossplane-tools`: v0.0.0-20240522174801-1ad3d4c87f21 → **needs upgrade to master**
- `upjet`: v1.4.1 → **needs upgrade to v2.1.0+**
- Go version: 1.23.0 ✅

### Directory Structure (Current)
```
provider-rancher/
├── apis/
│   ├── v1alpha1/           # Root API group with ProviderConfig
│   ├── app/                # Generated app resources
│   ├── auth/               # Generated auth resources
│   ├── k8s/                # Generated k8s resources
│   ├── rancher/            # Generated rancher resources
│   └── rbac/               # Generated rbac resources
├── config/
│   └── provider.go         # Single provider configuration
├── internal/controller/
│   ├── providerconfig/     # ProviderConfig controller
│   ├── app/
│   ├── auth/
│   ├── k8s/
│   ├── rancher/
│   └── rbac/
└── cmd/
    ├── generator/
    └── provider/
```

### Key Issues to Address
1. **StoreConfig (ESS)**: Currently present in `apis/v1alpha1/types.go` - must be removed (alpha feature, not supported in v2)
2. **Import paths**: All v1 imports need updating to v2
3. **Single provider config**: Need to support both cluster-scoped and namespaced MRs
4. **API Groups**: Need `.m.` variant for namespaced resources
5. **SafeStart capability**: Must be implemented for proper CRD availability checks

## Migration Steps

### Step 1: Update Dependencies ✋ **REQUIRES USER VALIDATION**

Update `go.mod`:
```go
require (
    github.com/crossplane/crossplane-runtime/v2 v2.0.0
    github.com/crossplane/crossplane-tools master
    github.com/crossplane/upjet/v2 v2.1.0
    // keep other dependencies
)
```

Run:
```bash
go mod tidy
```

### Step 2: Replace Import Paths ✋ **REQUIRES USER VALIDATION**

Search and replace across all `.go` files:
- `github.com/crossplane/crossplane-runtime/` → `github.com/crossplane/crossplane-runtime/v2/`
- `github.com/crossplane/upjet/` → `github.com/crossplane/upjet/v2/`

Affected files:
- `cmd/provider/main.go`
- `cmd/generator/main.go`
- `config/provider.go`
- `apis/v1alpha1/*.go`
- `internal/controller/**/*.go`

### Step 3: Remove External Secret Store (ESS) Support ✋ **REQUIRES USER VALIDATION**

**Files to modify:**
1. `apis/v1alpha1/types.go` - Remove `StoreConfig` type and related methods
2. `apis/v1alpha1/register.go` - Remove StoreConfig registration
3. `cmd/provider/main.go` - Remove ESS flags and logic (lines 57-144)

### Step 4: Restructure APIs Directory ✋ **REQUIRES USER VALIDATION**

**Actions:**
```bash
# 1. Move current apis to apis/cluster
mkdir -p apis/cluster
mv apis/app apis/cluster/
mv apis/auth apis/cluster/
mv apis/k8s apis/cluster/
mv apis/rancher apis/cluster/
mv apis/rbac apis/cluster/
mv apis/v1alpha1 apis/cluster/
mv apis/v1beta1 apis/cluster/ # if exists
mv apis/zz_register.go apis/cluster/

# 2. Keep generate.go at root
# apis/generate.go stays

# 3. Create namespaced directory
mkdir -p apis/namespaced/v1alpha1

# 4. Copy only root group to namespaced
cp -r apis/cluster/v1alpha1/* apis/namespaced/v1alpha1/
```

**Modify `apis/namespaced/v1alpha1/doc.go`:**
```go
// Package v1alpha1 contains the core resources of the Rancher namespaced provider.
// +kubebuilder:object:generate=true
// +groupName=rancher.m.contrib.crossplane.io
// +versionName=v1alpha1
package v1alpha1
```

**Modify `apis/namespaced/v1alpha1/types.go`:**
```go
package v1alpha1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
    xpv2 "github.com/crossplane/crossplane-runtime/v2/apis/common/v2"
)

const (
    Group   = "rancher.m.contrib.crossplane.io"
    Version = "v1alpha1"
)

// ProviderConfigSpec defines the desired state of a ProviderConfig.
type ProviderConfigSpec struct {
    // Credentials required to authenticate to this provider.
    Credentials ProviderCredentials `json:"credentials"`
}

type ProviderCredentials struct {
    // Source of the provider credentials.
    // +kubebuilder:validation:Enum=None;Secret;InjectedIdentity;Environment;Filesystem
    Source xpv1.CredentialsSource `json:"source"`

    xpv1.CommonCredentialSelectors `json:",inline"`
}

// ProviderConfigStatus defines the observed state of a ProviderConfig.
type ProviderConfigStatus struct {
    xpv1.ConditionedStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// A ProviderConfig configures the Rancher provider in a namespace.
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="SECRET-NAME",type="string",JSONPath=".spec.credentials.secretRef.name",priority=1
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:storageversion
type ProviderConfig struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   ProviderConfigSpec   `json:"spec"`
    Status ProviderConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProviderConfigList contains a list of ProviderConfig.
type ProviderConfigList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []ProviderConfig `json:"items"`
}

// +kubebuilder:object:root=true

// A ClusterProviderConfig configures the Rancher provider at cluster scope.
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="SECRET-NAME",type="string",JSONPath=".spec.credentials.secretRef.name",priority=1
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:storageversion
type ClusterProviderConfig struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   ProviderConfigSpec   `json:"spec"`
    Status ProviderConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterProviderConfigList contains a list of ClusterProviderConfig.
type ClusterProviderConfigList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []ClusterProviderConfig `json:"items"`
}

// +kubebuilder:object:root=true

// A ProviderConfigUsage indicates that a resource is using a ProviderConfig.
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="CONFIG-NAME",type="string",JSONPath=".providerConfigRef.name"
// +kubebuilder:printcolumn:name="RESOURCE-KIND",type="string",JSONPath=".resourceRef.kind"
// +kubebuilder:printcolumn:name="RESOURCE-NAME",type="string",JSONPath=".resourceRef.name"
// +kubebuilder:resource:scope=Namespaced,categories={crossplane}
type ProviderConfigUsage struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    xpv2.TypedProviderConfigUsage `json:",inline"`
}

// +kubebuilder:object:root=true

// ProviderConfigUsageList contains a list of ProviderConfigUsage
type ProviderConfigUsageList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []ProviderConfigUsage `json:"items"`
}
```

### Step 5: Restructure Config Directory ✋ **REQUIRES USER VALIDATION**

Current `config/` only has `provider.go` which is good - we'll keep the flat structure but add logic for both providers.

**Update `config/provider.go`:**
1. Add `GetProviderNamespaced()` function
2. Update root group for namespaced: `rancher.m.contrib.crossplane.io`
3. Both providers use same resource configurations (can be refactored later if needed)

### Step 6: Restructure Controllers ✋ **REQUIRES USER VALIDATION**

```bash
# 1. Move controllers to cluster
mkdir -p internal/controller/cluster
mv internal/controller/app internal/controller/cluster/
mv internal/controller/auth internal/controller/cluster/
mv internal/controller/k8s internal/controller/cluster/
mv internal/controller/rancher internal/controller/cluster/
mv internal/controller/rbac internal/controller/cluster/
mv internal/controller/providerconfig internal/controller/cluster/
mv internal/controller/zz_setup.go internal/controller/cluster/
mv internal/controller/doc.go internal/controller/cluster/

# 2. Create namespaced controllers
mkdir -p internal/controller/namespaced/providerconfig
```

**Create `internal/controller/cluster/providerconfig/config.go` SetupGated:**
```go
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
    o.Options.Gate.Register(func() {
        if err := Setup(mgr, o); err != nil {
            mgr.GetLogger().Error(err, "unable to setup reconciler", "gvk", v1alpha1.ProviderConfigGroupVersionKind.String())
        }
    }, v1alpha1.ProviderConfigGroupVersionKind, v1alpha1.ProviderConfigUsageGroupVersionKind)
    return nil
}
```

**Create namespaced providerconfig controller** similar structure.

### Step 7: Update Generator ✋ **REQUIRES USER VALIDATION**

**Modify `cmd/generator/main.go`:**
```go
import (
    "github.com/crossplane/upjet/v2/pkg/pipeline"
    "github.com/disaster37/provider-rancher/config"
)

func main() {
    // ... existing code ...
    pipeline.Run(config.GetProvider(), config.GetProviderNamespaced(), absRootDir)
}
```

### Step 8: Update Provider Main ✋ **REQUIRES USER VALIDATION**

**Major updates to `cmd/provider/main.go`:**
1. Import both cluster and namespaced APIs/controllers
2. Add k8s authv1 APIs to scheme
3. Remove ESS logic
4. Create separate options for cluster and namespaced
5. Implement SafeStart capability with CRD watching
6. Setup both controller sets

### Step 9: Update Makefile ✋ **REQUIRES USER VALIDATION**

Update versions in Makefile:
```makefile
KIND_VERSION = v0.30.0
UP_VERSION = v0.41.0
UP_CHANNEL = stable
CROSSPLANE_VERSION = 2.0.2
```

### Step 10: Update Package Metadata ✋ **REQUIRES USER VALIDATION**

**Modify `package/crossplane.yaml`:**
```yaml
apiVersion: meta.pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-rancher
  annotations:
    # ... existing annotations
spec:
  capabilities:
  - SafeStart
```

### Step 11: Generate Code ✋ **REQUIRES USER VALIDATION**

```bash
# Clean old generated code
rm -rf examples-generated/
find package/crds -name "*.yaml" -delete

# Regenerate
make generate
```

This will generate:
- `apis/cluster/` - cluster-scoped MR types  
- `apis/namespaced/` - namespace-scoped MR types
- `internal/controller/cluster/` - cluster-scoped controllers
- `internal/controller/namespaced/` - namespace-scoped controllers
- `package/crds/` - both cluster and namespaced CRDs

### Step 12: Create Example Resources ✋ **REQUIRES USER VALIDATION**

Create examples for new namespaced MRs in `examples/`:
- Use API group `*.rancher.m.contrib.crossplane.io`
- Add `metadata.namespace`
- Remove `namespace` from secret refs
- Update `spec.providerConfigRef` to typed reference

### Step 13: Testing ✋ **REQUIRES USER VALIDATION**

```bash
# Local deployment
make local-deploy

# Test cluster-scoped (legacy) MRs still work
kubectl apply -f examples/legacy/

# Test new namespaced MRs
kubectl apply -f examples/namespaced/
```

## Key Changes Summary

### API Changes
| Aspect | Before (v1) | After (v2) |
|--------|-------------|------------|
| Cluster API Group | `rancher.contrib.crossplane.io` | `rancher.contrib.crossplane.io` (unchanged) |
| Namespaced API Group | N/A | `rancher.m.contrib.crossplane.io` (.m. added) |
| ProviderConfig Types | 1 (cluster-scoped) | 3 (cluster ProviderConfig, namespaced ProviderConfig, ClusterProviderConfig) |
| MR Scope | Cluster only | Cluster + Namespaced |
| Secret References | Cross-namespace | Local namespace only (namespaced MRs) |
| ESS Support | Partial (alpha) | Removed |

### Breaking Changes
1. **StoreConfig removed** - External Secret Store is not supported
2. **publishConnectionDetailsTo removed** - Use `writeConnectionSecretToRef` only
3. **Secret references** - Namespaced MRs can only ref secrets in same namespace

### Backward Compatibility
✅ Existing cluster-scoped MRs continue to work  
✅ Can run on Crossplane v1  
✅ Namespaced MRs available but not composable in v1

## Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| Breaking existing resources | HIGH | Cluster-scoped MRs maintained, test thoroughly |
| Import path changes | MEDIUM | Systematic search/replace, compile check |
| Controller setup errors | HIGH | Careful implementation of SafeStart logic |
| CRD generation issues | MEDIUM | Clean regeneration, validate CRDs |
| Missing dependencies | LOW | Follow exact version requirements |

## Rollback Plan

If migration fails:
1. Revert `go.mod` changes
2. Restore original directory structure from backup
3. Run `go mod tidy && make generate`
4. Redeploy previous version

## Next Steps

I'm ready to help you execute this migration. Here's what I can do:

1. **Make structural changes** - Move directories, create new files
2. **Update import paths** - Search and replace across codebase  
3. **Modify configuration** - Update provider.go, main.go
4. **Generate new code** - Run make generate

**Each change will require your validation before proceeding to the next step.**

Would you like me to:
- Start with Step 1 (updating go.mod)?
- Create a backup branch first?
- Begin with a specific step you're most concerned about?
