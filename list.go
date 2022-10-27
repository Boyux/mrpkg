package mrpkg

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/Boyux/mrpkg/option"
	"net/url"
	"sort"
	"strings"
	"unsafe"
)

type ListIterator[T any] interface {
	Next() bool
	Value() T
}

func ToGoSlice[T any](iter ListIterator[T]) (goSlice []T) {
	goSlice = make([]T, 0, 20)
	for iter.Next() {
		goSlice = append(goSlice, iter.Value())
	}
	return goSlice
}

func Iter[T any](slice []T) *Slice[T] {
	var iter = Slice[T](slice)
	return &iter
}

type Slice[T any] []T

func (iter *Slice[T]) Next() bool {
	return len(*iter) > 0
}

func (iter *Slice[T]) Value() T {
	element := (*iter)[0]
	*iter = (*iter)[1:]
	return element
}

type Node[T any] struct {
	element list.Element
}

func (node *Node[T]) Value() T {
	return node.element.Value.(T)
}

func (node *Node[T]) Next() *Node[T] {
	return (*Node[T])(unsafe.Pointer(node.element.Next()))
}

func (node *Node[T]) Prev() *Node[T] {
	return (*Node[T])(unsafe.Pointer(node.element.Prev()))
}

func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{
		stdList: list.New(),
	}
}

type LinkedList[T any] struct {
	stdList *list.List
}

func (l *LinkedList[T]) Len() int {
	return l.stdList.Len()
}

func (l *LinkedList[T]) Front() *Node[T] {
	return (*Node[T])(unsafe.Pointer(l.stdList.Front()))
}

func (l *LinkedList[T]) Back() *Node[T] {
	return (*Node[T])(unsafe.Pointer(l.stdList.Back()))
}

func (l *LinkedList[T]) Remove(node *Node[T]) {
	l.stdList.Remove((*list.Element)(unsafe.Pointer(node)))
}

func (l *LinkedList[T]) PushFront(v T) {
	l.stdList.PushFront(v)
}

func (l *LinkedList[T]) PushBack(v T) {
	l.stdList.PushBack(v)
}

func (l *LinkedList[T]) InsertBefore(v T, mark *Node[T]) {
	l.stdList.InsertBefore(v, (*list.Element)(unsafe.Pointer(mark)))
}

func (l *LinkedList[T]) InsertAfter(v T, mark *Node[T]) {
	l.stdList.InsertAfter(v, (*list.Element)(unsafe.Pointer(mark)))
}

func (l *LinkedList[T]) MoveToFront(node *Node[T]) {
	l.stdList.MoveToFront((*list.Element)(unsafe.Pointer(node)))
}

func (l *LinkedList[T]) MoveToBack(node *Node[T]) {
	l.stdList.MoveToBack((*list.Element)(unsafe.Pointer(node)))
}

func (l *LinkedList[T]) MoveBefore(node, mark *Node[T]) {
	l.stdList.MoveBefore((*list.Element)(unsafe.Pointer(node)), (*list.Element)(unsafe.Pointer(mark)))
}

func (l *LinkedList[T]) MoveAfter(node, mark *Node[T]) {
	l.stdList.MoveAfter((*list.Element)(unsafe.Pointer(node)), (*list.Element)(unsafe.Pointer(mark)))
}

func (l *LinkedList[T]) PushBackList(other *LinkedList[T]) {
	l.stdList.PushBackList(other.stdList)
}

func (l *LinkedList[T]) PushFrontList(other *LinkedList[T]) {
	l.stdList.PushFrontList(other.stdList)
}

func (l *LinkedList[T]) sendToChan(ch chan<- *Node[T]) {
	for node := l.Front(); node != nil; node = node.Next() {
		ch <- node
	}
	close(ch)
}

func (l *LinkedList[T]) Iterator() <-chan *Node[T] {
	ch := make(chan *Node[T])
	go l.sendToChan(ch)
	return ch
}

type linkedListIterator[T any] struct {
	ch      <-chan *Node[T]
	hasNext bool
	current *Node[T]
}

func (iter *linkedListIterator[T]) Next() bool {
	return iter.hasNext
}

func (iter *linkedListIterator[T]) Value() T {
	node := iter.current

	if next, hasNext := <-iter.ch; hasNext {
		iter.hasNext = hasNext
		iter.current = next
	} else {
		iter.hasNext = false
		iter.current = nil
	}

	return node.Value()
}

func (l *LinkedList[T]) ListIterator() ListIterator[T] {
	ch := l.Iterator()
	current, hasNext := <-ch
	return &linkedListIterator[T]{
		ch:      ch,
		hasNext: hasNext,
		current: current,
	}
}

type Vector[T any] struct {
	mem []T
}

func (vector *Vector[T]) Len() int {
	return len(vector.mem)
}

func (vector *Vector[T]) Get(index int) option.Option[T] {
	if index < 0 || index >= vector.Len() {
		return option.None[T]()
	}

	return option.Some(vector.mem[index])
}

func (vector *Vector[T]) Push(element T) {
	vector.mem = append(vector.mem, element)
}

func (vector *Vector[T]) Append(list ListIterator[T]) {
	vector.mem = append(vector.mem, ToGoSlice(list)...)
}

func (vector *Vector[T]) SortBy(less func(l, r T) bool) {
	if vector.Len() == 0 {
		return
	}

	sort.Slice(vector.mem, func(i, j int) bool {
		return less(vector.mem[i], vector.mem[j])
	})
}

func (vector *Vector[T]) Insert(index int, element T) {
	originLen := vector.Len()
	if index > originLen {
		panic(fmt.Errorf("Vector.Insert: index %d out of range(%d)", index, originLen))
	} else if index < originLen {
		vector.mem = append(vector.mem, element)
		copy(vector.mem[index+1:], vector.mem[index:originLen])
		vector.mem[index] = element
	} else {
		vector.Push(element)
	}
}

func (vector *Vector[T]) sendToChan(ch chan<- T) {
	for i := 0; i < vector.Len(); i++ {
		ch <- vector.mem[i]
	}
	close(ch)
}

func (vector *Vector[T]) Iterator() <-chan T {
	ch := make(chan T)
	go vector.sendToChan(ch)
	return ch
}

type vectorIterator[T any] struct {
	mem []T
	idx int
}

func (iter *vectorIterator[T]) Next() bool {
	return iter.idx < len(iter.mem)
}

func (iter *vectorIterator[T]) Value() T {
	currentIdx := iter.idx
	iter.idx++
	return iter.mem[currentIdx]
}

func (vector *Vector[T]) ListIterator() ListIterator[T] {
	return &vectorIterator[T]{
		mem: vector.mem,
		idx: 0,
	}
}

type Strings []string

func (ss Strings) ToJSON() string {
	bytes, _ := json.Marshal(ss)
	return string(bytes)
}

func (ss Strings) Join(sep string) string {
	return strings.Join(ss, sep)
}

func (ss Strings) ToRawQuery(key string) string {
	query := make(url.Values)
	for _, s := range ss {
		query.Add(key, s)
	}
	return query.Encode()
}
