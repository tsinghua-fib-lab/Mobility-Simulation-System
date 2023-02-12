package crowd

import (
	"math"

	"git.fiblab.net/sim/simulet-go/entity"
	"git.fiblab.net/sim/simulet-go/utils/geometry"
)

// 行人在追踪的点的类型
type InterestType int32

const (
	InterestType_SLEEP InterestType = iota // 某个消失点
	InterestType_EXIT                      // 某个出口
)

// 室内行人状态
type Status int32

const (
	Status_IDLE         Status = iota // 运动中
	Status_PAUSE                      // 立定不动，做障碍物处理
	Status_REACH_TARGET               // 到达终点
)

// 运行时数据
type Runtime struct {
	Status       Status         // 状态
	Position     geometry.Point // 位置
	Velocity     geometry.Point // 速度
	Destination  geometry.Point // 目的地
	DesiresSpeed float64        // 期望速度
	Interest     InterestType   // 兴趣点类型，决定从crowd离开后去哪
	Person       entity.IPerson // 行人自身指针
}

// 每个行人的更新
func (r *Runtime) Update(positionMap map[int32]Snapshot, c *Crowd) {
	switch r.Status {
	case Status_IDLE:
		// 计算社会力，更新行人位置
		position := r.Position
		velocity := r.Velocity
		destination := r.Destination
		boundary := c.boundaryPoints
		// 记录当前行人所受合力
		force := geometry.Point{}

		// 计算期望移动目标提供的吸引力
		desiredDestination := destination.Sub(position)
		distance := math.Max(1e-5, desiredDestination.Length())
		desiredForce := desiredDestination.Scale(r.DesiresSpeed / distance).Sub(velocity)
		force.MoveVector(desiredForce, 2) // scale=2 表征吸引力的强度

		// 计算其他行人提供的排斥力
		for id, snapshot := range positionMap {
			if snapshot.Status == Status_REACH_TARGET || id == r.Person.ID() {
				continue
			}
			pos := snapshot.Position
			vec := snapshot.Velocity

			vr := pos.Sub(position)
			vv := vec.Sub(velocity)
			r := math.Max(1e-5, vr.Length())
			v := math.Max(1e-5, vv.Length())
			cos := geometry.Dot(vr, vv) / r / v
			// 之前手动调出的结果，若有需要可以将此超参的调整接口暴露出来
			const A, B, C, D, theta = 7.55, -3.0, .2, -.3, 56.0 / 180.0 * math.Pi
			f := A * math.Pow(math.E, B*r+C*cos+D*r*cos)
			direction := vr.Scale(-1 / r)
			direction = geometry.Point{
				X: math.Sin(math.Pi/2-theta)*direction.X - math.Sin(theta)*direction.Y,
				Y: math.Sin(theta)*direction.X + math.Sin(math.Pi/2-theta)*direction.Y,
			}
			force.MoveVector(direction, f)
		}

		// 计算AOI边界产生的排斥力
		// 找到障碍物边界上距离最近的点
		minDis := -1.0
		var closestPoint geometry.Point
		for i, j := 0, len(boundary)-1; i < len(boundary); i++ {
			p1, p2 := boundary[i], boundary[j]
			rn := geometry.Dot(position.Sub(p1), p2.Sub(p1))
			rd := geometry.SquareDistance(p2, p1)
			var r float64
			if rd < 1e-4 {
				// p1和p2基本重合，p1~p2不算做障碍物边界
				continue
			} else if rn < 0 {
				// p与线段p1~p2的最近点为p1
				r = 0
			} else if rn > rd {
				// p与线段p1~p2的最近点为p2
				r = 1
			} else {
				// p与线段p1~p2的最近点为二者中间某点
				r = rn / rd
			}
			cp := geometry.Blend(p1, p2, r)
			md := geometry.Distance(cp, position)
			if minDis < 0 || md < minDis {
				minDis = md
				closestPoint = cp
			}
			j = i + 1
		}
		minDis = math.Max(1e-5, minDis)
		// 计算该最近点产生的排斥力，若最近障碍点在门附近则不产生排斥力
		if minDis > 0 && !c.atAnyGate(closestPoint) {
			// 短程力
			const ShortA, ShortB, ShortTheta = 60.0, -3.0, 45 * math.Pi / 180
			f := ShortA * math.Pow(math.E, ShortB*minDis)
			closestPoint.MoveVector(position, -1)
			direction := closestPoint.Scale(-1 / minDis)
			vDesired := destination.Sub(position)
			// flag = 0: 行人正要远离障碍物，排斥力沿障碍物法向
			// flag = 1: 行人正要靠近障碍物，排斥力偏转，促使行人从左侧绕开障碍物
			// flag = -1: 行人正要靠近障碍物，排斥力偏转，促使行人从右侧绕开障碍物
			var flag int
			if geometry.Dot(vDesired, closestPoint) > 0 {
				flag = 1
			} else {
				flag = 0
			}
			if flag != 0 {
				if geometry.Cross(vDesired, closestPoint) > 0 {
					flag = -flag
				}
				direction = geometry.Point{
					X: math.Sin(math.Pi/2-ShortTheta)*direction.X + float64(flag)*math.Sin(ShortTheta)*direction.Y,
					Y: -float64(flag)*math.Sin(ShortTheta)*direction.X + math.Sin(math.Pi/2-ShortTheta)*direction.Y,
				}
			}
			force.MoveVector(direction, f)

			// 长程力
			const LongA, LongB, LongTheta = .1, -.5, 75 * math.Pi / 180
			f = LongA * math.Pow(math.E, LongB*minDis)
			direction = closestPoint.Scale(-1 / minDis)
			if flag != 0 {
				direction = geometry.Point{
					X: math.Sin(math.Pi/2-LongTheta)*direction.X + float64(flag)*math.Sin(LongTheta)*direction.Y,
					Y: -float64(flag)*math.Sin(LongTheta)*direction.X + math.Sin(math.Pi/2-LongTheta)*direction.Y,
				}
			}
			force.MoveVector(direction, f)
		}

		acc := force
		velocity.MoveVector(acc, STEP_INTERVAL_IN_DOOR)
		speed := velocity.Length()
		if speed > 1.2*r.DesiresSpeed {
			velocity = velocity.Scale(1.2 * r.DesiresSpeed / speed)
		}
		position.MoveVector(velocity, STEP_INTERVAL_IN_DOOR)
		r.Position = position
		r.Velocity = velocity

		// 以小概率转为原地站定
		if r.Person.Generator().PTrue(.01) {
			r.Status = Status_PAUSE
		}

		// 到达终点
		if r.DistanceToEnd() < 5.0 {
			r.Status = Status_REACH_TARGET
		}
	case Status_PAUSE:
		r.Velocity = geometry.Point{}
		// 以小概率转为开始运动
		const PAUSE_TO_IDLE_P = .09
		if r.Person.Generator().PTrue(PAUSE_TO_IDLE_P) {
			r.Status = Status_IDLE
		}
	case Status_REACH_TARGET:
	}
}

// 获取到终点的距离
func (r *Runtime) DistanceToEnd() float64 {
	return geometry.Distance(r.Position, r.Destination)
}
