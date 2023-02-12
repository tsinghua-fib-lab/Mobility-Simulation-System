package trafficlight

import (
	"fmt"

	"git.fiblab.net/sim/simulet-go/entity"
	lightv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/traffic_light/v2"
	"git.fiblab.net/sim/simulet-go/utils"
)

type localTlRuntime struct {
	tl              *lightv2.TrafficLight
	tlStep          int32
	tlRemainingTime float64
}

type localTrafficLight struct {
	JunctionID int32                            // 所属junction ID
	lanes      []entity.ILaneTrafficLightSetter // 车道数据

	timeBeforeChange [][]float64     // 下一次信号灯变化时间（相位切换时不一定所有的信号灯都变）
	snapshot         localTlRuntime  // snapshot，用于保存输出的数据
	runtime          localTlRuntime  // 运行时数据
	buffer           *localTlRuntime // 数据buffer，用于交互式接口写入(optional)
	ok               bool            // 信号灯状态，true为开启，false为关闭
	okBuffer         bool            // 信号灯状态buffer，用于交互式接口写入
}

func NewLocalTrafficLight(junctionID int32, lanes []entity.ILaneTrafficLightSetter) *localTrafficLight {
	return &localTrafficLight{
		JunctionID:       junctionID,
		lanes:            lanes,
		timeBeforeChange: make([][]float64, 0),
		runtime:          localTlRuntime{},
		ok:               true,
		okBuffer:         true,
	}
}

func (l *localTrafficLight) Prepare() {
	// 更新信号灯状态
	l.ok = l.okBuffer
	// 写入snapshot
	l.snapshot = l.runtime
	// 写入lane中数据
	if l.snapshot.tl == nil || !l.ok {
		for _, lane := range l.lanes {
			lane.SetLightState(lightv2.LightState_LIGHT_STATE_GREEN)
			lane.SetLightRemainingTime(utils.INF)
		}
	} else {
		p := l.snapshot.tl.Phases[l.snapshot.tlStep]
		for i, lane := range l.lanes {
			lane.SetLightState(p.States[i])
			lane.SetLightRemainingTime(
				l.snapshot.tlRemainingTime + l.timeBeforeChange[i][l.snapshot.tlStep],
			)
		}
	}
}

func (l *localTrafficLight) Update(stepInterval float64) {
	if l.buffer != nil {
		l.runtime = *l.buffer
		l.buffer = nil
		// 初始化步骤
		if l.runtime.tl != nil {
			n := len(l.runtime.tl.Phases)
			j := len(l.lanes)
			for i := 0; i < j; i++ {
				time := make([]float64, n)
				allTheSame := true
				lastState := l.runtime.tl.Phases[n-1].States[i]
				for k := n - 2; k >= 0; k-- {
					state := l.runtime.tl.Phases[k+1].States[i]
					if state == lastState {
						time[k] = time[k+1] + l.runtime.tl.Phases[k+1].Duration
					} else {
						allTheSame = false
					}
					lastState = state
				}
				if allTheSame {
					time = make([]float64, n)
					for idx := 0; idx < n; idx++ {
						time[idx] = utils.INF
					}
				} else {
					t0 := time[0] + l.runtime.tl.Phases[0].Duration
					for k := n - 1; k >= 0; k-- {
						if lastState != l.runtime.tl.Phases[k].States[i] {
							break
						}
						time[k] += t0
					}
				}
				l.timeBeforeChange = append(l.timeBeforeChange, time)
			}
		}
	}
	if l.runtime.tl == nil || !l.ok {
		return
	}
	l.runtime.tlRemainingTime -= stepInterval
	// 切换相位
	if l.runtime.tlRemainingTime <= 0 {
		for {
			l.runtime.tlStep = (l.runtime.tlStep + 1) % int32(len(l.runtime.tl.Phases))
			l.runtime.tlRemainingTime += l.runtime.tl.Phases[l.runtime.tlStep].Duration
			if l.runtime.tlRemainingTime > 0 {
				break
			}
		}
	}
}

func (l *localTrafficLight) Set(tl *lightv2.TrafficLight) error {
	if tl.JunctionId != l.JunctionID {
		return fmt.Errorf("set junction %d with wrong traffic light id %d", l.JunctionID, tl.JunctionId)
	}
	if l.lanes == nil {
		return fmt.Errorf("no lane data in junction %d", l.JunctionID)
	}
	if tl.Phases == nil {
		return fmt.Errorf("set with empty traffic light")
	}
	for _, p := range tl.Phases {
		if len(p.States) != len(l.lanes) {
			return fmt.Errorf("number of lanes %d and traffic light states %d does not match", len(l.lanes), len(p.States))
		}
	}

	phaseIndex := l.JunctionID % int32(len(tl.Phases))
	l.buffer = &localTlRuntime{
		tl: tl, tlStep: phaseIndex, tlRemainingTime: tl.Phases[phaseIndex].Duration,
	}
	return nil
}
