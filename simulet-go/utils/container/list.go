package container

// sorted double linked list

type ListNode[T any, E any] struct {
	parent     *List[T, E]
	prev, next *ListNode[T, E]
	Key        float64
	Value      T
	Extra      E
}

func (n *ListNode[T, E]) Prev() *ListNode[T, E] {
	return n.prev
}

func (n *ListNode[T, E]) Next() *ListNode[T, E] {
	return n.next
}

func (n *ListNode[T, E]) ValueOrDefault(defaultValue T) T {
	if n == nil {
		return defaultValue
	}
	return n.Value
}

func (n *ListNode[T, E]) InsertBefore(add *ListNode[T, E]) {
	add.parent = n.parent
	add.next = n
	add.prev = n.prev
	n.prev = add
	if add.prev != nil {
		add.prev.next = add
	} else {
		add.parent.head = add
	}
	n.parent.length++
}

func (n *ListNode[T, E]) InsertAfter(add *ListNode[T, E]) {
	add.parent = n.parent
	add.prev = n
	add.next = n.next
	n.next = add
	if add.next != nil {
		add.next.prev = add
	} else {
		add.parent.tail = add
	}
	n.parent.length++
}

type List[T any, E any] struct {
	head, tail *ListNode[T, E]
	length     int
}

func NewList[T any, E any]() *List[T, E] {
	l := &List[T, E]{head: nil, tail: nil, length: 0}
	return l
}

func (l *List[T, E]) Keys() []float64 {
	keys := make([]float64, 0)
	for node := l.head; node != nil; node = node.next {
		keys = append(keys, node.Key)
	}
	return keys
}

func (l *List[T, E]) Values() []T {
	values := make([]T, 0)
	for node := l.head; node != nil; node = node.next {
		values = append(values, node.Value)
	}
	return values
}

func (l *List[T, E]) Length() int {
	return l.length
}

func (l *List[T, E]) PushFront(add *ListNode[T, E]) {
	add.parent = l
	add.next = nil
	add.prev = nil
	if l.head == nil {
		l.head = add
		l.tail = add
		l.length++
	} else {
		// length++在InsertBefore中处理
		l.head.InsertBefore(add)
		l.head = add
	}
}

func (l *List[T, E]) PushBack(add *ListNode[T, E]) {
	add.parent = l
	add.next = nil
	add.prev = nil
	if l.tail == nil {
		l.head = add
		l.tail = add
		l.length++
	} else {
		// length++在InsertAfter中处理
		l.tail.InsertAfter(add)
		l.tail = add
	}
}

func (l *List[T, E]) Remove(node *ListNode[T, E]) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		l.head = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		l.tail = node.prev
	}
	node.prev = nil
	node.next = nil
	node.parent = nil
	l.length--
}

func (l *List[T, E]) First() *ListNode[T, E] {
	return l.head
}

func (l *List[T, E]) Last() *ListNode[T, E] {
	return l.tail
}

// 移除逆序节点
func (l *List[T, E]) PopUnsorted() []*ListNode[T, E] {
	unsorted := make([]*ListNode[T, E], 0)
	for node := l.head; node != nil; {
		next := node.next
		if node.prev != nil && node.prev.Key > node.Key {
			l.Remove(node)
			unsorted = append(unsorted, node)
		}
		node = next
	}
	return unsorted
}

func (l *List[T, E]) Merge(adds []*ListNode[T, E]) {
	// 1. sort array (可优化)
	for i := 0; i < len(adds)-1; i++ {
		for j := i + 1; j < len(adds); j++ {
			if adds[i].Key > adds[j].Key {
				adds[i], adds[j] = adds[j], adds[i]
			}
		}
	}
	// 2. merge sort
	node := l.head
	for _, add := range adds {
		for node != nil && node.Key < add.Key {
			node = node.next
		}
		if node != nil {
			node.InsertBefore(add)
		} else {
			l.PushBack(add)
		}
	}
}
