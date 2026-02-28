package enterprise

// Registry is the extension point for enterprise features.
// The open-source build uses the no-op implementation below.
// Enterprise builds replace this via build tags.
type Registry struct{}

// NewRegistry returns an initialised Registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// IsEnabled returns whether an enterprise feature flag is active.
// Always returns false in the OSS build.
func (r *Registry) IsEnabled(_ string) bool {
	return false
}
