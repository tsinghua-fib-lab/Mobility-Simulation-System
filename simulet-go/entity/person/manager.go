package person

import (
	"sync"

	agentv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/agent/v2"
	"git.fiblab.net/sim/simulet-go/utils/clock"
	"git.fiblab.net/sim/simulet-go/utils/container"
	"git.fiblab.net/sim/simulet-go/utils/parallel"
)

const (
	BLOCK_SIZE = 512
)

var (
	Manager = NewManager()
)

type PersonManager struct {
	data map[int32]*Person

	// 有计算、输出需求的person
	vehicles    *container.ActiveContainer[*Person]
	pedestrians *container.ActiveContainer[*Person]
	crowds      *container.ActiveContainer[*Person] // 只用于输出

	personInserted      []*Person // 新加入的人
	personInsertedMutex sync.Mutex
	nextPersonID        int32
}

func NewManager() *PersonManager {
	return &PersonManager{
		data:                make(map[int32]*Person),
		vehicles:            container.NewActiveData[*Person](),
		pedestrians:         container.NewActiveData[*Person](),
		crowds:              container.NewActiveData[*Person](),
		personInserted:      make([]*Person, 0),
		personInsertedMutex: sync.Mutex{},
		nextPersonID:        10000000,
	}
}

func (m *PersonManager) Init(agents []*agentv2.Agent) {
	for _, agent := range agents {
		m.data[agent.Id] = NewPerson(agent, m)
	}
}

func (m *PersonManager) Data() map[int32]*Person {
	return m.data
}

func (m *PersonManager) Prepare() {
	// 新人加入
	for _, newP := range m.personInserted {
		if _, ok := m.data[newP.ID()]; ok {
			log.Panic("Person: same id between new person and existed person")
		}
		m.data[newP.ID()] = newP
	}
	m.personInserted = []*Person{}

	// data prepare
	// 最好不要并行处理，因为共用index，如果一个人同时从车辆中删去又加入行人，可能有问题
	m.vehicles.Prepare()
	m.pedestrians.Prepare()
	m.crowds.Prepare()

	// active person prepare
	var wg sync.WaitGroup
	parallel.GoForWithWaitGroup(m.vehicles.Data(), 2048, &wg, func(p *Person) { p.Prepare() })
	parallel.GoForWithWaitGroup(m.pedestrians.Data(), 2048, &wg, func(p *Person) { p.Prepare() })
	wg.Wait()
}

func (m *PersonManager) Update(stepInterval float64) {
	// 性能测试
	if clock.Step%100 == 0 {
		log.Errorf("STEP: %d #vehicles: %d #pedestrians: %d", clock.Step, len(m.vehicles.Data()), len(m.pedestrians.Data()))
	}

	var wg sync.WaitGroup
	parallel.GoForWithWaitGroup(m.vehicles.Data(), 1024, &wg, func(person *Person) { person.Update(stepInterval) })
	parallel.GoForWithWaitGroup(m.pedestrians.Data(), 2048, &wg, func(person *Person) { person.Update(stepInterval) })
	wg.Wait()
}
