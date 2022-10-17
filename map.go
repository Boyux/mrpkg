package mrpkg

import (
	"golang.org/x/exp/constraints"
	"sort"
	"sync"
)

type OrderedMapKey interface {
	comparable
	constraints.Ordered
}

type MapIterator[K any, V any] interface {
	Next() bool
	Value() (K, V)
}

func ToGoMap[K comparable, V any](iter MapIterator[K, V]) (goMap map[K]V) {
	goMap = make(map[K]V, 20)
	for iter.Next() {
		k, v := iter.Value()
		goMap[k] = v
	}
	return goMap
}

func OrderedMap[K OrderedMapKey, V any](m map[K]V) MapIterator[K, V] {
	keys := MapKeys(m)
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return &orderedMap[K, V]{
		Idx:  0,
		Keys: keys,
		Map:  m,
	}
}

type orderedMap[K comparable, V any] struct {
	Idx  int
	Keys []K
	Map  map[K]V
}

func (o *orderedMap[K, V]) Next() bool {
	return o.Idx < len(o.Keys)
}

func (o *orderedMap[K, V]) Value() (K, V) {
	k := o.Keys[o.Idx]
	v := o.Map[k]
	o.Idx++
	return k, v
}

type Entry[K any, V any] struct {
	Key   K
	Value V
	Map   *ConcurrentMap[K, V]
}

// ConcurrentMap is a lock free concurrent map
type ConcurrentMap[K any, V any] struct {
	DefaultFunc func(K) V
	syncMap     sync.Map
}

func (m *ConcurrentMap[K, V]) Len() (n int) {
	m.syncMap.Range(func(_, _ any) bool {
		n++
		return true
	})
	return n
}

func (m *ConcurrentMap[K, V]) Set(key K, val V) {
	m.syncMap.Store(key, val)
}

func (m *ConcurrentMap[K, V]) Del(key K) {
	m.syncMap.Delete(key)
}

func (m *ConcurrentMap[K, V]) Get(key K) (val V, ok bool) {
	var v any
	v, ok = m.syncMap.Load(key)
	if !ok {
		return
	}
	return v.(V), true
}

func (m *ConcurrentMap[K, V]) GetOrDefault(key K) (val V) {
	var defaultFunc = m.DefaultFunc
	if defaultFunc == nil {
		defaultFunc = func(K) V {
			return New[V]()
		}
	}
	v, _ := m.syncMap.LoadOrStore(key, defaultFunc(key))
	return v.(V)
}

func (m *ConcurrentMap[K, V]) sendToChan(ch chan<- *Entry[K, V]) {
	m.syncMap.Range(func(key, value any) bool {
		ch <- &Entry[K, V]{
			Key:   key.(K),
			Value: value.(V),
			Map:   m,
		}
		return true
	})
	close(ch)
}

func (m *ConcurrentMap[K, V]) Iterator() <-chan *Entry[K, V] {
	ch := make(chan *Entry[K, V])
	go m.sendToChan(ch)
	return ch
}

func (m *ConcurrentMap[K, V]) Keys() (keys []K) {
	keys = make([]K, 0, m.Len())
	for entry := range m.Iterator() {
		keys = append(keys, entry.Key)
	}
	return keys
}

func (m *ConcurrentMap[K, V]) Values() (values []V) {
	values = make([]V, 0, m.Len())
	for entry := range m.Iterator() {
		values = append(values, entry.Value)
	}
	return values
}

type concurrentMapIterator[K any, V any] struct {
	ch      <-chan *Entry[K, V]
	hasNext bool
	current *Entry[K, V]
}

func (iter *concurrentMapIterator[K, V]) Next() bool {
	return iter.hasNext
}

func (iter *concurrentMapIterator[K, V]) Value() (K, V) {
	entry := iter.current

	if next, hasNext := <-iter.ch; hasNext {
		iter.hasNext = hasNext
		iter.current = next
	} else {
		iter.hasNext = false
		iter.current = nil
	}

	return entry.Key, entry.Value
}

func (m *ConcurrentMap[K, V]) MapIterator() MapIterator[K, V] {
	ch := m.Iterator()
	current, hasNext := <-ch
	return &concurrentMapIterator[K, V]{
		ch:      ch,
		hasNext: hasNext,
		current: current,
	}
}

func NewConcurrentSet[T any](hashFunc func(T) any) *ConcurrentSet[T] {
	set := new(ConcurrentSet[T])
	set.HashFunc = hashFunc
	return set
}

type ConcurrentSet[T any] struct {
	concurrentMap ConcurrentMap[any, T]
	HashFunc      func(T) any
}

func (set *ConcurrentSet[T]) hash(x T) any {
	if set.HashFunc == nil {
		return x
	}
	return set.HashFunc(x)
}

func (set *ConcurrentSet[T]) Len() int {
	return set.concurrentMap.Len()
}

func (set *ConcurrentSet[T]) Add(element T) {
	set.concurrentMap.Set(set.hash(element), element)
}

func (set *ConcurrentSet[T]) Del(element T) {
	set.concurrentMap.Del(set.hash(element))
}

func (set *ConcurrentSet[T]) BatchAdd(iter ListIterator[T]) {
	for iter.Next() {
		set.Add(iter.Value())
	}
}

func (set *ConcurrentSet[T]) BatchDel(iter ListIterator[T]) {
	for iter.Next() {
		set.Del(iter.Value())
	}
}

func (set *ConcurrentSet[T]) Union(other *ConcurrentSet[T]) (target *ConcurrentSet[T]) {
	target = new(ConcurrentSet[T])
	target.BatchAdd(set.ListIterator())
	target.BatchAdd(other.ListIterator())
	return target
}

func (set *ConcurrentSet[T]) Intersection(other *ConcurrentSet[T]) (target *ConcurrentSet[T]) {
	target = new(ConcurrentSet[T])
	target.BatchAdd(set.ListIterator())
	target.BatchDel(set.Difference(other).ListIterator())
	return target
}

func (set *ConcurrentSet[T]) Difference(other *ConcurrentSet[T]) (target *ConcurrentSet[T]) {
	target = new(ConcurrentSet[T])
	target.BatchAdd(set.ListIterator())
	target.BatchDel(other.ListIterator())
	return target
}

func (set *ConcurrentSet[T]) SymmetricDifference(other *ConcurrentSet[T]) (target *ConcurrentSet[T]) {
	target = new(ConcurrentSet[T])
	target.BatchAdd(set.Union(other).ListIterator())
	target.BatchDel(set.Intersection(other).ListIterator())
	return target
}

func (set *ConcurrentSet[T]) Contains(element T) bool {
	_, loaded := set.concurrentMap.Get(set.hash(element))
	return loaded
}

func (set *ConcurrentSet[T]) sendToChan(ch chan<- T) {
	for entry := range set.concurrentMap.Iterator() {
		ch <- entry.Value
	}
	close(ch)
}

func (set *ConcurrentSet[T]) Iterator() <-chan T {
	ch := make(chan T)
	go set.sendToChan(ch)
	return ch
}

type concurrentSetIterator[T any] struct {
	mapIterator MapIterator[any, T]
}

func (iter *concurrentSetIterator[T]) Next() bool {
	return iter.mapIterator.Next()
}

func (iter *concurrentSetIterator[T]) Value() T {
	_, element := iter.mapIterator.Value()
	return element
}

func (set *ConcurrentSet[T]) ListIterator() ListIterator[T] {
	return &concurrentSetIterator[T]{
		mapIterator: set.concurrentMap.MapIterator(),
	}
}
