package mrpkg

import (
	"github.com/Boyux/mrpkg/option"
	"golang.org/x/exp/constraints"
	"sync/atomic"
)

type Cell[T any] struct {
	inner atomic.Value
}

func (cell *Cell[T]) Store(value T) {
	cell.inner.Store(value)
}

func (cell *Cell[T]) Load() (value option.Option[T]) {
	if v := cell.inner.Load(); v != nil {
		return option.Some(v.(T))
	}
	return option.None[T]()
}

func (cell *Cell[T]) CompareAndSwap(old, new T) (swapped bool) {
	return cell.inner.CompareAndSwap(old, new)
}

func AtomicAdd[T constraints.Integer](cell *Cell[T], value T) T {
	if v := cell.Load(); option.IsNull(v) {
		if cell.inner.CompareAndSwap(nil, value) {
			return value
		}
	}

	for {
		oldValue := cell.Load().Value()
		newValue := oldValue + value
		if cell.CompareAndSwap(oldValue, newValue) {
			return newValue
		}
	}
}
