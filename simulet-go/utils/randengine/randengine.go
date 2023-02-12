package randengine

import (
	"golang.org/x/exp/rand"
	"sync"
)

type Engine struct {
	*rand.Rand
	mtx sync.Mutex
}

func New(seed uint64) *Engine {
	e := &Engine{Rand: rand.New(rand.NewSource(seed))}
	return e
}

// 按给定概率分布生成 [0, n) 中的随机数，非多线程安全
func (e *Engine) DiscreteDistribution(weight []float64) int32 {
	random := .0
	for _, w := range weight {
		random += w
	}
	random *= e.Float64()
	sum := .0
	for i, w := range weight {
		sum += w
		if sum > random {
			return int32(i)
		}
	}
	return int32(len(weight))
}

// 以p概率给出true
func (e *Engine) PTrue(p float64) bool {
	return e.Float64() < p
}

// 以p概率给出true，多线程安全
func (e *Engine) PTrueSafe(p float64) bool {
	e.mtx.Lock()
	defer e.mtx.Unlock()
	return e.Float64() < p
}

func (e *Engine) IntnSafe(n int) int {
	e.mtx.Lock()
	defer e.mtx.Unlock()
	return e.Intn(n)
}
