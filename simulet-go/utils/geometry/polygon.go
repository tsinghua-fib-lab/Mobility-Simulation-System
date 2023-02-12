package geometry

// TODO: 可以考虑自行实现，这个库用的人不多
import (
	"github.com/ctessum/geom"
	"github.com/samber/lo"
)

func GetPolygonCentroid(polygon []Point) Point {
	if len(polygon) == 1 {
		return polygon[0]
	}
	geoPath := geom.Path(lo.Map(polygon, func(point Point, _ int) geom.Point {
		return *geom.NewPoint(point.X, point.Y)
	}))
	geoPolygon := geom.Polygon([]geom.Path{geoPath})
	c := geoPolygon.Centroid()
	return Point{X: c.X, Y: c.Y}
}
