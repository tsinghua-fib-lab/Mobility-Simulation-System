package geometry

import (
	"math"
)

// get the "s" of the closest point in line to point
func getClosestLineSToPoint(lineStart, lineEnd, point Point) float64 {
	// vector a = point - line_start
	ax, ay := point.X-lineStart.X, point.Y-lineStart.Y
	// vector b = line_end - line_start
	bx, by := lineEnd.X-lineStart.X, lineEnd.Y-lineStart.Y
	// s = clamp(dot(a, unit(b)), 0, |b|) = clamp(dot(a, b) / |b|, 0, |b|)
	bl := math.Sqrt(bx*bx + by*by)
	return math.Max(0, math.Min(bl, (ax*bx+ay*by)/bl))
}

func getDistanceToLine(lineStart, lineEnd, point Point) float64 {
	// vector a = point - line_start
	ax, ay := point.X-lineStart.X, point.Y-lineStart.Y
	// vector b = line_end - line_start
	bx, by := lineEnd.X-lineStart.X, lineEnd.Y-lineStart.Y
	// ratio = dot(a, unit(b)) / |b| = dot(a, b) / |b|²
	// ratio = clamp(ratio, 0, 1)
	b2 := bx*bx + by*by
	ratio := math.Max(0, math.Min(1, ax*bx+ay*by)/b2)
	// vector d = vector a - vector b * ratio
	dx := ax - bx*ratio
	dy := ay - by*ratio
	return math.Sqrt(dx*dx + dy*dy)
}

// compute the length (a.k.a s) of each node
// the first is 0 and the last is the length of the polyline
func GetPolylineLengths(line []Point) []float64 {
	result := make([]float64, 0, len(line))
	s := 0.0
	result = append(result, s)
	p1 := line[0]
	for _, p2 := range line[1:] {
		s += math.Sqrt((p2.X-p1.X)*(p2.X-p1.X) + (p2.Y-p1.Y)*(p2.Y-p1.Y))
		result = append(result, s)
		p1 = p2
	}
	return result
}

func GetPolylineDirections(line []Point) []float64 {
	result := make([]float64, 0, len(line)-1)
	for i, p1 := range line[:len(line)-1] {
		p2 := line[i+1]
		result = append(result, math.Atan2(p2.Y-p1.Y, p2.X-p1.X))
	}
	return result
}

// project point (x, y) to line and return "s" at the polyline
// line_lengths: GetPolylineLengths
func GetClosestPolylineSToPoint(line []Point, lineLengths []float64, point Point) float64 {
	index := 0
	dMin := getDistanceToLine(line[0], line[1], point)
	for i, j, size := 1, 2, len(line); j < size; i, j = i+1, j+1 {
		d := getDistanceToLine(line[i], line[j], point)
		if d < dMin {
			index = i
			dMin = d
		}
	}
	s := lineLengths[index] + getClosestLineSToPoint(line[index], line[index+1], point)
	return math.Max(lineLengths[index], math.Min(lineLengths[index+1], s))
}
