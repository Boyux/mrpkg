package mrpkg

import "container/heap"

type Queue[T any] interface {
	Len() int
	Push(T)
	Pop() T
}

type Less[T any] interface {
	Lt(T) bool
}

type tinyHeap[T Less[T]] []T

func (hp *tinyHeap[T]) Len() int {
	return len(*hp)
}

func (hp *tinyHeap[T]) Less(i, j int) bool {
	return (*hp)[i].Lt((*hp)[j])
}

func (hp *tinyHeap[T]) Swap(i, j int) {
	(*hp)[i], (*hp)[j] = (*hp)[j], (*hp)[i]
}

func (hp *tinyHeap[T]) Push(x any) {
	*hp = append(*hp, x.(T))
}

func (hp *tinyHeap[T]) Pop() any {
	x := (*hp)[hp.Len()-1]
	*hp = (*hp)[:hp.Len()-1]
	return x
}

type Heap[T Less[T]] struct {
	inner tinyHeap[T]
}

func (hp *Heap[T]) Len() int {
	return hp.inner.Len()
}

func (hp *Heap[T]) Push(x T) {
	heap.Push(&hp.inner, x)
}

func (hp *Heap[T]) Pop() T {
	return heap.Pop(&hp.inner).(T)
}
