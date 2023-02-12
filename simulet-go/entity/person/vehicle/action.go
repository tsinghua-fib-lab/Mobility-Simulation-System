package vehicle

type Action struct {
	Acc              float64
	EnableLaneChange bool
	LaneChangeLength float64
}

// 采用取最小的方式设置加速度
func (a *Action) UpdateByMinAcc(others ...Action) {
	for _, o := range others {
		if o.Acc < a.Acc {
			a.Acc = o.Acc
		}
	}
}

func (a *Action) StartLaneChange(lcLength float64) {
	a.EnableLaneChange = true
	a.LaneChangeLength = lcLength
}
