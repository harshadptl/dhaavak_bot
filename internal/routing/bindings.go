package routing

// BindingStore holds the configured agent-channel bindings.
type BindingStore struct {
	bindings     []Binding
	defaultAgent string
}

// NewBindingStore creates a store from config bindings and a default agent ID.
func NewBindingStore(bindings []Binding, defaultAgent string) *BindingStore {
	return &BindingStore{
		bindings:     bindings,
		defaultAgent: defaultAgent,
	}
}

// Bindings returns all configured bindings.
func (bs *BindingStore) Bindings() []Binding {
	return bs.bindings
}

// DefaultAgent returns the fallback agent ID.
func (bs *BindingStore) DefaultAgent() string {
	return bs.defaultAgent
}
