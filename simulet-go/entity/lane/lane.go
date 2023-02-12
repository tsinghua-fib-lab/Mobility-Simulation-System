package lane

import (
	"sort"
	"sync"

	"git.fiblab.net/sim/simulet-go/entity"
	"git.fiblab.net/sim/simulet-go/utils"
	"git.fiblab.net/sim/simulet-go/utils/container"

	geov2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/geo/v2"
	mapv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/map/v2"
	lightv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/traffic_light/v2"
	"git.fiblab.net/sim/simulet-go/utils/geometry"
	"github.com/samber/lo"
)

type Lane struct {
	id int32

	// 初始化临时变量

	initPredecessors []*mapv2.LaneConnection
	initSuccessors   []*mapv2.LaneConnection
	initLeftLaneIDs  []int32
	initRightLaneIDs []int32

	maxSpeed       float64                     // 当前道路限速
	parentJunction entity.IJunction            // 所在路口
	predecessors   map[int32]entity.Connection // 前驱车道映射表
	successors     map[int32]entity.Connection // 后继车道映射表
	sideLanes      [2][]entity.ILane           // 左/右侧车道（按距离从近到远排序）
	aois           map[int32]entity.IAoi       // AOI映射表
	addAoiMutex    sync.Mutex                  // aois读写互斥锁
	lineLengths    []float64                   // 中心线折线点对应的的长度列表
	length         float64                     // 以中心线的长度为车道长度
	width          float64                     // 车道宽度
	lineDirections []float64                   // 中心线折线段每一段的方向（atan2）
	line           []geometry.Point            // 转成Point的中心线折线

	pedestrians laneList[entity.IPedestrian, struct{}]
	vehicles    laneList[entity.IVehicle, entity.VehicleSideLink]

	lightState              lightv2.LightState // 车道信号灯状态
	lightStateRemainingTime float64            // 车道信号灯下一次切换时间
}

func NewLane(base *mapv2.Lane) *Lane {
	l := &Lane{
		id:                      base.Id,
		initPredecessors:        base.Predecessors,
		initSuccessors:          base.Successors,
		initLeftLaneIDs:         base.LeftLaneIds,
		initRightLaneIDs:        base.RightLaneIds,
		maxSpeed:                base.MaxSpeed,
		predecessors:            make(map[int32]entity.Connection),
		successors:              make(map[int32]entity.Connection),
		sideLanes:               [2][]entity.ILane{},
		aois:                    make(map[int32]entity.IAoi),
		addAoiMutex:             sync.Mutex{},
		lineLengths:             make([]float64, 0),
		lineDirections:          make([]float64, 0),
		line:                    make([]geometry.Point, 0),
		width:                   base.Width,
		pedestrians:             newLaneList[entity.IPedestrian, struct{}](),
		vehicles:                newLaneList[entity.IVehicle, entity.VehicleSideLink](),
		lightState:              lightv2.LightState_LIGHT_STATE_GREEN,
		lightStateRemainingTime: utils.INF,
	}
	l.line = lo.Map(base.CenterLine.Nodes, func(node *geov2.XYPosition, _ int) geometry.Point {
		return geometry.NewPointFromPb(node)
	})
	l.lineLengths = geometry.GetPolylineLengths(l.line)
	l.length = l.lineLengths[len(l.lineLengths)-1]
	l.lineDirections = geometry.GetPolylineDirections(l.line)
	return l
}

func (l *Lane) InitLanes(laneManager *LaneManager) {
	for _, conn := range l.initPredecessors {
		if lane, err := laneManager.Get(conn.Id); err != nil {
			log.Panic(err)
		} else {
			l.predecessors[conn.Id] = entity.Connection{Lane: lane, Type: conn.Type}
		}
	}
	for _, conn := range l.initSuccessors {
		if lane, err := laneManager.Get(conn.Id); err != nil {
			log.Panic(err)
		} else {
			l.successors[conn.Id] = entity.Connection{Lane: lane, Type: conn.Type}
		}
	}
	for _, id := range l.initLeftLaneIDs {
		if lane, err := laneManager.Get(id); err != nil {
			log.Panic(err)
		} else {
			l.sideLanes[entity.LEFT] = append(l.sideLanes[entity.LEFT], lane)
		}
	}
	for _, id := range l.initRightLaneIDs {
		if lane, err := laneManager.Get(id); err != nil {
			log.Panic(err)
		} else {
			l.sideLanes[entity.RIGHT] = append(l.sideLanes[entity.RIGHT], lane)
		}
	}
	l.initPredecessors = nil
	l.initSuccessors = nil
	l.initLeftLaneIDs = nil
	l.initRightLaneIDs = nil
}

func (l *Lane) Prepare() {
	// 维护本车道链表
	l.pedestrians.prepare()
	l.vehicles.prepare()
}

func (l *Lane) Prepare2() {
	// 等待相邻车道完成主链构建并进行支链构建
	for which := range []int{entity.LEFT, entity.RIGHT} {
		thisSideLanes := l.sideLanes[which]
		if len(thisSideLanes) > 0 {
			neighborLane := thisSideLanes[0]
			// 根据邻居车道链表构建本车道链表支链
			nList := neighborLane.Vehicles()
			var nBack *container.ListNode[entity.IVehicle, entity.VehicleSideLink] = nil
			nFront := nList.First()
			if nFront == nil {
				// 隔壁车道没车，不需要任何处理
				return
			}
			for node := l.vehicles.list.First(); node != nil; node = node.Next() {
				// 找到第一个位置大于node的邻居车道上的车
				// 则nFront是第一个位置大于等于node的车，nBack是第一个位置小于node的车
				// 该算法能处理nFront和nBack为nil的情况
				for nFront != nil && l.ProjectFromLane(neighborLane, nFront.Key) < node.Key {
					nBack = nFront
					nFront = nFront.Next()
				}
				node.Extra.Links[which][entity.BEFORE] = nBack
				node.Extra.Links[which][entity.AFTER] = nFront
			}
		}
	}
}

// 数据初始化

func (l *Lane) SetParentJunctionWhenInit(parent entity.IJunction) {
	l.parentJunction = parent
}

func (l *Lane) AddAoiWhenInit(aoi entity.IAoi) {
	l.addAoiMutex.Lock()
	l.aois[aoi.ID()] = aoi
	l.addAoiMutex.Unlock()
}

// 静态数据

func (l *Lane) ID() int32 {
	return l.id
}

func (l *Lane) Length() float64 {
	return l.length
}

func (l *Lane) Width() float64 {
	return l.width
}

func (l *Lane) Aois() map[int32]entity.IAoi {
	return l.aois
}

func (l *Lane) Successors() map[int32]entity.Connection {
	return l.successors
}

func (l *Lane) Predecessors() map[int32]entity.Connection {
	return l.predecessors
}

func (l *Lane) InJunction() bool {
	return l.parentJunction != nil
}

func (l *Lane) FirstLeftLane() entity.ILane {
	if len(l.sideLanes[entity.LEFT]) == 0 {
		return nil
	} else {
		return l.sideLanes[entity.LEFT][0]
	}
}

func (l *Lane) FirstRightLane() entity.ILane {
	if len(l.sideLanes[entity.RIGHT]) == 0 {
		return nil
	} else {
		return l.sideLanes[entity.RIGHT][0]
	}
}

// 信号灯

func (l *Lane) LightState() lightv2.LightState {
	return l.lightState
}

func (l *Lane) SetLightState(state lightv2.LightState) {
	l.lightState = state
}

func (l *Lane) LightStateRemainingTime() float64 {
	return l.lightStateRemainingTime
}

func (l *Lane) SetLightRemainingTime(time float64) {
	l.lightStateRemainingTime = time
}

// 路况

func (l *Lane) MaxSpeed() float64 {
	return l.maxSpeed
}

// 人车更新相关函数

func (l *Lane) Vehicles() *container.List[entity.IVehicle, entity.VehicleSideLink] {
	return l.vehicles.list
}

func (l *Lane) Pedestrians() *container.List[entity.IPedestrian, struct{}] {
	return l.pedestrians.list
}

func (l *Lane) ReportPedestrianAdded(node *container.ListNode[entity.IPedestrian, struct{}]) {
	l.pedestrians.reportAdded(node)
}

func (l *Lane) ReportPedestrianRemoved(node *container.ListNode[entity.IPedestrian, struct{}]) {
	l.pedestrians.reportRemoved(node)
}

func (l *Lane) ReportVehicleAdded(node *container.ListNode[entity.IVehicle, entity.VehicleSideLink]) {
	l.vehicles.reportAdded(node)
}

func (l *Lane) ReportVehicleRemoved(node *container.ListNode[entity.IVehicle, entity.VehicleSideLink]) {
	l.vehicles.reportRemoved(node)
}

// 获取第一辆车
func (l *Lane) GetFirstVehicle() entity.IVehicle {
	return l.vehicles.list.First().ValueOrDefault(nil)
}

// 对同一道路内的车道按比例"投影"
func (l *Lane) ProjectFromLane(other entity.ILane, otherS float64) float64 {
	return lo.Clamp(otherS/other.Length()*l.length, 0, l.length)
}

// 根据本车道s坐标计算切向角度
func (l *Lane) GetDirectionByS(s float64) (direction float64) {
	if s < l.lineLengths[0] || s > l.lineLengths[len(l.lineLengths)-1] {
		log.Warnf("get direction with s %v out of range{%v,%v}",
			s, l.lineLengths[0], l.lineLengths[len(l.lineLengths)-1])
		s = lo.Clamp(s, l.lineLengths[0], l.lineLengths[len(l.lineLengths)-1])
	}
	if i := sort.SearchFloat64s(l.lineLengths, s); i == 0 {
		direction = l.lineDirections[0]
	} else {
		direction = l.lineDirections[i-1]
	}
	return
}

// 将当前车道s坐标转换为xy坐标
func (l *Lane) GetPositionByS(s float64) (pos geometry.Point) {
	if s < l.lineLengths[0] || s > l.lineLengths[len(l.lineLengths)-1] {
		log.Warnf("get position with s %v out of range{%v,%v}",
			s, l.lineLengths[0], l.lineLengths[len(l.lineLengths)-1])
		s = lo.Clamp(s, l.lineLengths[0], l.lineLengths[len(l.lineLengths)-1])
	}
	if i := sort.SearchFloat64s(l.lineLengths, s); i == 0 {
		pos = l.line[0]
	} else {
		sHigh, sSlow := l.lineLengths[i], l.lineLengths[i-1]
		pos = geometry.Blend(l.line[i-1], l.line[i], (s-sSlow)/(sHigh-sSlow))
	}
	return
}
