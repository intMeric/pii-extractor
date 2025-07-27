package extractors

import (
	"fmt"
	"sync"
)

// Registry manages available PII extractors
type Registry struct {
	extractors map[string]PiiExtractor
	mu         sync.RWMutex
}

// NewRegistry creates a new extractor registry
func NewRegistry() *Registry {
	return &Registry{
		extractors: make(map[string]PiiExtractor),
	}
}

// Register adds an extractor to the registry
func (r *Registry) Register(name string, extractor PiiExtractor) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if extractor == nil {
		return fmt.Errorf("cannot register nil extractor")
	}
	
	r.extractors[name] = extractor
	return nil
}

// Get retrieves an extractor by name
func (r *Registry) Get(name string) (PiiExtractor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	extractor, exists := r.extractors[name]
	if !exists {
		return nil, fmt.Errorf("extractor '%s' not found", name)
	}
	
	return extractor, nil
}

// List returns all registered extractor names
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	names := make([]string, 0, len(r.extractors))
	for name := range r.extractors {
		names = append(names, name)
	}
	
	return names
}

// GetByMethod returns all extractors that use the specified method
func (r *Registry) GetByMethod(method ExtractionMethod) []PiiExtractor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var extractors []PiiExtractor
	for _, extractor := range r.extractors {
		if extractor.GetMethod() == method {
			extractors = append(extractors, extractor)
		}
	}
	
	return extractors
}

// Default global registry
var defaultRegistry = NewRegistry()

// Register adds an extractor to the default registry
func Register(name string, extractor PiiExtractor) error {
	return defaultRegistry.Register(name, extractor)
}

// Get retrieves an extractor from the default registry
func Get(name string) (PiiExtractor, error) {
	return defaultRegistry.Get(name)
}

// List returns all extractor names from the default registry
func List() []string {
	return defaultRegistry.List()
}

// GetByMethod returns extractors by method from the default registry
func GetByMethod(method ExtractionMethod) []PiiExtractor {
	return defaultRegistry.GetByMethod(method)
}