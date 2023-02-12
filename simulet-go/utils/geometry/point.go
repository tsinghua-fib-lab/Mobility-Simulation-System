package geometry

import (
	"math"

	geov2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/geo/v2"
)

// 参考：https://github.com/golang/geo/blob/master/r2/rect.go
// https://stackoverflow.com/questions/27775376/value-receiver-vs-pointer-receiver
// If the receiver is a small array or struct that is naturally a value type (for instance, something like the time.Time type), with no mutable fields and no pointers, or is just a simple basic type such as int or string, a value receiver makes sense.
// A value receiver can reduce the amount of garbage that can be generated; if a value is passed to a value method, an on-stack copy can be used instead of allocating on the heap. (The compiler tries to be smart about avoiding this allocation, but it can't always succeed.) Don't choose a value receiver type for this reason without profiling first.

// 二维向量
type Point struct {
	X, Y float64
}

func NewPointFromPb(pb *geov2.XYPosition) Point {
	return Point{X: pb.X, Y: pb.Y}
}

// 向量+
func (a Point) Add(b Point) Point { return Point{a.X + b.X, a.Y + b.Y} }

// 向量-
func (a Point) Sub(b Point) Point { return Point{a.X - b.X, a.Y - b.Y} }

// 数乘
func (p Point) Scale(k float64) Point { return Point{k * p.X, k * p.Y} }

// 向量长度平方
func (p Point) SquareLength() float64 { return p.X*p.X + p.Y*p.Y }

// 向量长度
func (p Point) Length() float64 { return math.Sqrt(p.SquareLength()) }

// 向量角度
func (p Point) Angle() float64 { return math.Atan2(p.Y, p.X) }

// 指定方向移动指定距离
func (p *Point) MoveDirection(direction, distance float64) {
	p.X += distance * math.Cos(direction)
	p.Y += distance * math.Sin(direction)
}

// 指定向量移动指定距离
func (p *Point) MoveVector(v Point, scale float64) {
	p.X += v.X * scale
	p.Y += v.Y * scale
}

// return start*(1-k) + b*k, 0<=k<=1
func Blend(start, end Point, k float64) Point {
	return Point{start.X*(1-k) + end.X*k, start.Y*(1-k) + end.Y*k}
}

func dcmp(x float64) int32 {
	if (math.Abs(x)) < 1e-6 {
		return 0
	} else if x < 0 {
		return -1
	} else {
		return 1
	}
}

// 用计算几何中的射线法判断点是否在任意多边形内部（不含边）
// 取水平向右的射线，计算与多边形边的交点数，根据结果的奇偶性判断是否在内部
// 请保证positions各点顺序给出，且第一点与最后一点相同
func (p Point) InPolygon(positions []Point) bool {
	flag := false
	for i, j := 0, len(positions)-2; i < len(positions)-1; i++ {
		p1, p2 := positions[i], positions[j]
		onSegment := dcmp(Cross(p1.Sub(p), p2.Sub(p))) == 0 &&
			dcmp(Dot(p1.Sub(p), p2.Sub(p))) <= 0
		if onSegment {
			return false
		}
		if ((dcmp(p1.Y-p.Y) > 0) != (dcmp(p2.Y-p.Y) > 0)) &&
			dcmp((p.X-p1.X)-(p1.X-p2.X)*(p.Y-p1.Y)/(p1.Y-p2.Y)) < 0 {
			flag = !flag
		}
		j = i
	}
	return flag
}

// 点乘
func Dot(a, b Point) float64 { return a.X*b.Y + a.Y*b.Y }

// 叉乘
func Cross(a, b Point) float64 { return a.X*b.Y - b.X*a.Y }

func SquareDistance(a, b Point) float64 { return (a.X-b.X)*(a.X-b.X) + (a.Y-b.Y)*(a.Y-b.Y) }

func Distance(a, b Point) float64 { return math.Sqrt(SquareDistance(a, b)) }
