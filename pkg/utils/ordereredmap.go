package utils

type OrderedMap struct {
	keys   []string
	values map[string]interface{}
}

// NewOrderedMap creates a new ordered map.
func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		keys:   []string{},
		values: make(map[string]interface{}),
	}
}

// Set adds a key-value pair to the map, preserving insertion order.
func (om *OrderedMap) Set(key string, value interface{}) {
	// Only add key to the keys slice if it's not already present
	if _, exists := om.values[key]; !exists {
		om.keys = append(om.keys, key)
	}
	om.values[key] = value
}

// Get retrieves a value by key.
func (om *OrderedMap) Get(key string) (interface{}, bool) {
	val, exists := om.values[key]
	return val, exists
}

// Remove deletes a key-value pair from the map and removes the key from the order.
func (om *OrderedMap) Remove(key string) {
	// Remove key from the keys slice
	for i, k := range om.keys {
		if k == key {
			om.keys = append(om.keys[:i], om.keys[i+1:]...)
			break
		}
	}
	delete(om.values, key)
}

// Keys returns the keys in the order they were added.
func (om *OrderedMap) Keys() []string {
	return om.keys
}

// Values returns the values in the order of their corresponding keys.
func (om *OrderedMap) Values() []interface{} {
	var vals []interface{}
	for _, key := range om.keys {
		vals = append(vals, om.values[key])
	}
	return vals
}
