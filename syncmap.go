package winc

import "sync"

type SyncMap[K comparable, V any] struct {
	_map sync.Map
}

func (m *SyncMap[K, V]) Delete(key K) {
	m._map.Delete(key)
}

func (m *SyncMap[K, V]) Load(key K) (V, bool) {
	val, ok := m._map.Load(key)
	if val == nil {
		return *new(V), ok
	}

	return val.(V), ok
}

func (m *SyncMap[K, V]) LoadAndDelete(key K) (V, bool) {
	val, ok := m._map.LoadAndDelete(key)
	if val == nil {
		return *new(V), ok
	}

	return val.(V), ok
}

func (m *SyncMap[K, V]) LoadOrStore(key K, value V) (V, bool) {
	val, ok := m._map.LoadOrStore(key, value)
	if val == nil {
		return *new(V), ok
	}

	return val.(V), ok
}

func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m._map.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}

func (m *SyncMap[K, V]) Store(key K, value V) {
	m._map.Store(key, value)
}
