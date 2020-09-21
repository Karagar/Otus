package hw04_lru_cache //nolint:golint,stylecheck

import "sync"

type List interface {
	Len() int                          // длина списка
	Front() *listItem                  // первый listItem
	Back() *listItem                   // последний listItem
	PushFront(v interface{}) *listItem // добавить значение в начало
	PushBack(v interface{}) *listItem  // добавить значение в конец
	Remove(i *listItem)                // удалить элемент
	MoveToFront(i *listItem)           // переместить элемент в начало
}

type listItem struct {
	Value    interface{}
	CacheKey string
	Prev     *listItem
	Next     *listItem
}

type list struct {
	length   int
	first    *listItem
	last     *listItem
	listLock sync.Mutex
}

// Len (length) of double-linked list.
func (l *list) Len() int {
	return l.length
}

// First listItem of double-linked list.
func (l *list) Front() *listItem {
	return l.first
}

// Last listItem of double-linked list.
func (l *list) Back() *listItem {
	return l.last
}

// PushFront added one more listItem element to the front of double-linked list.
func (l *list) PushFront(v interface{}) *listItem {
	item := listItem{
		Value: v,
		Prev:  nil,
		Next:  l.first,
	}
	l.listLock.Lock()
	if l.length != 0 {
		l.first.Prev = &item
	} else {
		l.last = &item
	}
	l.first = &item
	l.length++
	l.listLock.Unlock()
	return &item
}

// PushBack added one more listItem element to the back of double-linked list.
func (l *list) PushBack(v interface{}) *listItem {
	item := listItem{
		Value: v,
		Prev:  l.last,
		Next:  nil,
	}
	l.listLock.Lock()
	if l.length != 0 {
		l.last.Next = &item
	} else {
		l.first = &item
	}
	l.last = &item
	l.length++
	l.listLock.Unlock()
	return &item
}

// Remove removed one listItem element from double-linked list.
func (l *list) Remove(i *listItem) {
	l.listLock.Lock()
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if l.last == i {
		l.last = i.Prev
	}
	if l.first == i {
		l.first = i.Next
	}
	l.length--
	l.listLock.Unlock()
}

// MoveToFront.
func (l *list) MoveToFront(i *listItem) {
	l.Remove(i)
	l.listLock.Lock()
	i.Prev = nil
	i.Next = l.first
	if l.length != 0 {
		l.first.Prev = i
	} else {
		l.last = i
	}
	l.first = i
	l.length++
	l.listLock.Unlock()
}

func NewList() List {
	return &list{}
}
