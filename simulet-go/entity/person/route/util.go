package route

import "git.fiblab.net/sim/simulet-go/entity"

// 导航起点，lane+s/aoi二选一
type RouteStartPosition struct {
	Lane entity.ILane
	S    float64
	Aoi  entity.IAoi
}
