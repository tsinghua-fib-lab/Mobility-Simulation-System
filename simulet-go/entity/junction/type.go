package junction

import (
	lightv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/traffic_light/v2"
)

// 依赖倒置，表达junction对信号灯实现的接口需求

type ITrafficLight interface {
	Prepare()                           // 处理各种写入buffer，将信控结果写入到lane中
	Update(stepInterval float64)        // 更新信控结果
	Set(tl *lightv2.TrafficLight) error // 修改信控程序
}
