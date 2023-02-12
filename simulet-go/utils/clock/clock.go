package clock

import "git.fiblab.net/sim/simulet-go/utils/config"

// 全局时钟信息

var (
	STEP_INTERVAL float64 = 0 // 每个模拟步时间间隔
	START_STEP    int32   = 0 // 起始步
	END_STEP      int32   = 0 // 结束步，模拟区间[START, END)

	GlobalTime float64 = 0 // 当前时间
	Step       int32   = 0 // 当前步
)

func Init(c config.Config) {
	STEP_INTERVAL = c.Control.Step.Interval
	START_STEP = c.Control.Step.Start
	END_STEP = c.Control.Step.Start + c.Control.Step.Total

	Step = START_STEP
	GlobalTime = float64(Step) * STEP_INTERVAL
}
