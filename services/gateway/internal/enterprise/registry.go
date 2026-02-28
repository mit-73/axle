package enterprise

// Registry is the extension point for enterprise features in the Gateway.
// The open-source build uses the no-op implementation below.
type Registry struct{}

// NewRegistry returns an initialised Registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// IsEnabled returns false in the OSS build.
func (r *Registry) IsEnabled(_ string) bool {
	return false
}
