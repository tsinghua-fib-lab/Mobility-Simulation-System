package person

import (
	"git.fiblab.net/sim/simulet-go/entity"
)

// person子模块
type submodule interface {
	Update(stepInterval float64)
	Prepare()

	// Person从子模块更新自身snapshot

	FetchBaseSnapshotForPerson() (entity.BaseRuntime, entity.BaseRuntimeOnRoad)

	// 生命周期检查与辅助函数

	GetEndByPerson() (endAoi entity.IAoi, endLane entity.ILane, isEnd bool) // 该子模块生命周期结束，提供终点信息
}
