package nationalgrid

import (
	"fmt"

	geos "github.com/twpayne/go-geos"
)

type Quadrant string

const (
	NE Quadrant = "NE"
	NW Quadrant = "NW"
	SE Quadrant = "SE"
	SW Quadrant = "SW"
)

func (l Quadrant) String() string {
	if l == NE {
		return "NE"
	}
	if l == NW {
		return "NW"
	}
	if l == SE {
		return "SE"
	}
	if l == SW {
		return "SW"
	}
	return ""
}

type GridRef struct {
	Square    string
	SubSquare string
	Quadrant  Quadrant
}

type GridSquare struct {
	Geom *geos.Geom
	MinX float64
	MinY float64
}

type Bounds struct {
	Xmin float64
	Xmax float64
	Ymin float64
	Ymax float64
}

func (b Bounds) ToPolygon() string {
	tl := []float64{
		b.Xmin,
		b.Ymin,
	}

	tr := []float64{
		b.Xmax,
		b.Ymin,
	}

	bl := []float64{
		b.Xmin,
		b.Ymax,
	}

	br := []float64{
		b.Xmax,
		b.Ymax,
	}

	return fmt.Sprintf(
		"POLYGON((%v %v, %v %v, %v %v, %v %v, %v %v))",
		tl[0], tl[1], tr[0], tr[1], br[0], br[1], bl[0], br[1], tl[0], tl[1],
	)
}
