package container

import "sync"

type IActiveElement interface {
	Index() int
	SetIndex(index int)
}

type ActiveElement struct {
	index int
}

func (a *ActiveElement) Index() int {
	return a.index
}

func (a *ActiveElement) SetIndex(index int) {
	a.index = index
}

type ActiveContainer[T IActiveElement] struct {
	data        []T
	add         []T
	remove      []T
	addMutex    sync.Mutex
	removeMutex sync.Mutex
}

func NewActiveData[T IActiveElement]() *ActiveContainer[T] {
	return &ActiveContainer[T]{
		data:   make([]T, 0),
		add:    make([]T, 0),
		remove: make([]T, 0),
	}
}

func (d *ActiveContainer[T]) Data() []T {
	return d.data
}

func (d *ActiveContainer[T]) MarkAsAdded(value T) {
	d.addMutex.Lock()
	defer d.addMutex.Unlock()
	d.add = append(d.add, value)
}

func (d *ActiveContainer[T]) MarkAsRemoved(value T) {
	d.removeMutex.Lock()
	defer d.removeMutex.Unlock()
	d.remove = append(d.remove, value)
}

func (d *ActiveContainer[T]) Prepare() {
	// 增 >= 删
	if len(d.add) >= len(d.remove) {
		addI := 0
		for _, removeV := range d.remove {
			addV := d.add[addI]
			ind := removeV.Index()
			addV.SetIndex(ind)
			d.data[ind] = addV
			addI++
		}
		for ; addI < len(d.add); addI++ {
			addV := d.add[addI]
			addV.SetIndex(len(d.data))
			d.data = append(d.data, addV)
		}
	} else {
		// 删 > 增
		removeI := 0
		for _, addV := range d.add {
			removeV := d.remove[removeI]
			ind := removeV.Index()
			addV.SetIndex(ind)
			d.data[ind] = addV
			removeI++
		}
		for ; removeI < len(d.remove); removeI++ {
			// 与最后一个做交换
			swapV := d.data[len(d.data)-1]
			ind := d.remove[removeI].Index()
			swapV.SetIndex(ind)
			d.data[ind] = swapV
			d.data = d.data[:len(d.data)-1]
		}
	}

	d.add = []T{}
	d.remove = []T{}
}
