# Upjet v2 Migration - Progress Report & Next Steps

## ✅ Completed Steps (4/17) - 24% Complete

### Step 1: Dependencies Updated ✅
- `crossplane-runtime/v2 v2.0.0` ✅
- `upjet/v2 v2.2.0` ✅  
- K8s packages v0.33.0 ✅
- All dependencies resolving correctly ✅

### Step 2: Import Paths Migrated ✅
- All 99+ `.go` files updated to v2 paths ✅
- Both crossplane-runtime and upjet imports fixed ✅
- Verified with `go mod tidy` ✅

### Step 3: StoreConfig/ESS Removed ✅  
- [`apis/v1alpha1/types.go`](apis/v1alpha1/types.go:1) - removed StoreConfig types ✅
- [`apis/v1alpha1/register.go`](apis/v1alpha1/register.go:1) - removed registration ✅  
- [`cmd/provider/main.go`](cmd/provider/main.go:1) - removed ESS flags & logic ✅

### Step 4: Directory Restructure ✅
- Created `apis/cluster` and `apis/namespaced` structure ✅
- Existing APIs already in `apis/cluster/` ✅
- Copied v1alpha1 to `apis/namespaced/v1alpha1/` ✅

**Files in namespaced/v1alpha1:**
- `doc.go` - needs API group update
- `register.go` - needs API group update  
- `types.go` - needs ClusterProviderConfig addition
- `zz_generated.deepcopy.go` - will be regenerated

---

## 🔄 Critical Next Steps

### Step 5: Update API Group Markers (REQUIRED)

**File:** `apis/namespaced/v1alpha1/doc.go`
Change:
```go
// +groupName=rancher.contrib.crossplane.io
```
To:
```go
// +groupName=rancher.m.contrib.crossplane.io
```

**File:** `apis/namespaced/v1alpha1/register.go`
Update const:
```go
const (
	Group   = "rancher.m.contrib.crossplane.io"  // ADD .m. here
	Version = "v1alpha1"
)
```

### Step 6: Add ClusterProviderConfig Type (REQUIRED)

**File:** `apis/namespaced/v1alpha1/types.go`

Need to add:
1. Import `xpv2` for TypedProviderConfigUsage:
```go
import (
	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	xpv2 "github.com/crossplane/crossplane-runtime/v2/apis/common/v2"
)
```

2. Update ProviderConfigUsage:
```go
type ProviderConfigUsage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	xpv2.TypedProviderConfigUsage `json:",inline"`  // CHANGED from xpv1.ProviderConfigUsage
}
```

3. Add ClusterProviderConfig (duplicate ProviderConfig but cluster-scoped):
```go
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:storageversion
type ClusterProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderConfigSpec   `json:"spec"`
	Status ProviderConfigStatus `json:"status,omitempty"`
}

type ClusterProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterProviderConfig `json:"items"`
}
```

4. Update ProviderConfig scope marker:
```go
// +kubebuilder:resource:scope=Namespaced  // CHANGED from Cluster
type ProviderConfig struct {
	...
}
```

5. Register ClusterProviderConfig in register.go:
```go
func init() {
	SchemeBuilder.Register(&ProviderConfig{}, &ProviderConfig List{})
	SchemeBuilder.Register(&ClusterProviderConfig{}, &ClusterProviderConfigList{})
	SchemeBuilder.Register(&ProviderConfigUsage{}, &ProviderConfigUsageList{})
}

// Add metadata vars
var (
	ClusterProviderConfigKind             = reflect.TypeOf(ClusterProviderConfig{}).Name()
	ClusterProviderConfigGroupKind        = schema.GroupKind{Group: Group, Kind: ClusterProviderConfigKind}.String()
	ClusterProviderConfigKindAPIVersion   = ClusterProviderConfigKind + "." + SchemeGroupVersion.String()
	ClusterProviderConfigGroupVersionKind = SchemeGroupVersion.WithKind(ClusterProviderConfigKind)
)
```

### Remaining Steps (7-17)

**Step 7:** Config directory restructure (separate cluster/namespaced configs)
**Step 8:** Add `GetProviderNamespaced()` to config/provider.go
**Step 9:** Update cmd/generator/main.go to call both providers  
**Step 10:** Restructure controllers (cluster + namespaced)
**Step 11:** Add SetupGated functions to controllers
**Step 12:** Major rewrite of cmd/provider/main.go for SafeStart
**Step 13:** Update Makefile versions
**Step 14:** Add SafeStart to package/crossplane.yaml
**Step 15:** **Regenerate all code** (fixes compilation errors)
**Step 16:** Create namespaced examples
**Step 17:** Test & validate

---

## ⚠️ Current State

**Compilation Status:** ❌ Multiple errors (expected)
- Types still reference removed Store Config
- Generated deepcopy methods out of sync
- Feature flags need v2 API updates
- **These will be fixed in Step 15 (regeneration)**

**Directory Structure:** ✅ Correct
```
apis/
├── generate.go          # Stays at root
├── cluster/             # ✅ Cluster-scoped MRs
│   ├── v1alpha1/
│   ├── app/
│   ├── auth/
│   ├── k8s/
│   ├── rancher/
│   └── rbac/
└── namespaced/          # ✅ Namespace-scoped MRs  
    └── v1alpha1/        # ✅ Needs API group + ClusterProviderConfig updates
```

---

## 📋 Quick Action Items

**Immediate (Steps 5-6):**
1. Update API group in namespaced/v1alpha1 files (.m. suffix)
2. Add ClusterProviderConfig type
3. Update ProviderConfigUsage to use xpv2.TypedProviderConfigUsage

**Next Phase (Steps 7-9):**
1. Organize config/ directory
2. Add namespaced provider function
3. Update generator

**Final Phase (Steps 10-17):**
1. Controller refactoring
2. SafeStart implementation
3. Code regeneration (will fix all compilation errors)
4. Testing

---

## 📊 Estimated Time Remaining

- Immediate steps (5-6): 30-45 min
- Next phase (7-9): 45-60 min
- Final phase (10-17): 2-3 hours
- **Total remaining: ~4 hours**

---

## 🔧 Ready to Continue

The migration is 24% complete (4/17 steps). The foundation is solid - dependencies updated, imports fixed, ESS removed, and directory structure ready.

**Next:** Update API group markers and add ClusterProviderConfig type (Steps 5-6), then proceed systematically through remaining steps.

See [`MIGRATION_PLAN.md`](MIGRATION_PLAN.md:1) for detailed implementation guide for each step.
