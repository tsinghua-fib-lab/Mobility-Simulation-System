package lane

import (
	"sync"

	"git.fiblab.net/sim/simulet-go/utils/container"
)

type laneList[T any, E any] struct {
	list              *container.List[T, E]
	addBuffer         []*container.ListNode[T, E]
	addBufferMutex    sync.Mutex
	removeBuffer      []*container.ListNode[T, E]
	removeBufferMutex sync.Mutex
}

func newLaneList[T any, E any]() laneList[T, E] {
	return laneList[T, E]{
		list:              container.NewList[T, E](),
		addBuffer:         make([]*container.ListNode[T, E], 0),
		addBufferMutex:    sync.Mutex{},
		removeBuffer:      make([]*container.ListNode[T, E], 0),
		removeBufferMutex: sync.Mutex{},
	}
}

func (l *laneList[T, E]) prepare() {
	for _, v := range l.removeBuffer {
		l.list.Remove(v)
	}
	unsorted := l.list.PopUnsorted()
	l.list.Merge(append(l.addBuffer, unsorted...))
	l.removeBuffer = l.removeBuffer[:0]
	l.addBuffer = l.addBuffer[:0]
}

func (l *laneList[T, E]) reportAdded(node *container.ListNode[T, E]) {
	l.addBufferMutex.Lock()
	l.addBuffer = append(l.addBuffer, node)
	l.addBufferMutex.Unlock()
}

func (l *laneList[T, E]) reportRemoved(node *container.ListNode[T, E]) {
	l.removeBufferMutex.Lock()
	l.removeBuffer = append(l.removeBuffer, node)
	l.removeBufferMutex.Unlock()
}
