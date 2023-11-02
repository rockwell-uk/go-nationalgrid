package nationalgrid

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/rockwell-uk/go-geos-draw/geom"
	geos "github.com/twpayne/go-geos"
)

var (
	gridSquares     = make(map[string]*geos.Bounds)
	subSquareCoords = make(map[string]map[int]*geos.Bounds)
)

func init() {
	// We need a context here
	gctx := geos.NewContext()

	for key, gridSquare := range NationalGridSquares {
		tileX := gridSquare[0] * SquareSize
		tileY := gridSquare[1] * SquareSize

		boundingBox := Bounds{
			Xmin: tileX,
			Xmax: tileX + SquareSize,
			Ymin: tileY,
			Ymax: tileY + SquareSize,
		}

		gridSquarePoly := boundingBox.ToPolygon()
		gridSquareGeom, _ := gctx.NewGeomFromWKT(gridSquarePoly)

		gridSquares[key] = gridSquareGeom.Bounds()

		subSquareCoords[key] = make(map[int]*geos.Bounds)
	}

	for key, gridSquare := range gridSquares {
		for i := 0; i <= 99; i++ {
			bits := fmt.Sprintf("%02d", i)
			addX, _ := strconv.Atoi(string(bits[0]))
			addY, _ := strconv.Atoi(string(bits[1]))

			xAdj := gridSquare.MinX + (float64(addX) * SubSquareSize)
			yAdj := gridSquare.MinY + (float64(addY) * SubSquareSize)

			boundingBox := Bounds{
				Xmin: xAdj,
				Xmax: xAdj + SubSquareSize,
				Ymin: yAdj,
				Ymax: yAdj + SubSquareSize,
			}

			subSquarePolygon := boundingBox.ToPolygon()
			subSquareGeom, _ := gctx.NewGeomFromWKT(subSquarePolygon)

			subSquareCoords[key][i] = subSquareGeom.Bounds()
		}
	}
}

// determine which squares / subsquares a geometry from a shapefile is within.
func GetSubSquares(g *geos.Bounds) map[string][]int {
	subSquares := make(map[string][]int)

	for key, gridSquare := range gridSquares {
		k := strings.ToLower(key)

		if intersects(g, gridSquare) {
			for i, subSquare := range subSquareCoords[key] {
				if intersects(g, subSquare) {
					subSquares[k] = append(subSquares[k], i)
					sort.Ints(subSquares[k])
				}
			}
		}
	}

	return subSquares
}

func intersects(item *geos.Bounds, target *geos.Bounds) bool {
	return item.Intersects(target)
}

func GetGridEastNorth(ref string) (float64, float64, error) {
	var east, north float64

	gridRef, err := ParseGridRef(ref)
	if err != nil {
		return east, north, err
	}

	return getGridCoordCenter(gridRef)
}

func getGridCoordCenter(gridRef GridRef) (float64, float64, error) {
	var x, y float64
	var g *geos.Geom
	var gridCoords []float64
	var err error

	var squareBlx, squareBly float64
	var subSquareBlx, subSquareBly float64
	var quadrantBlx, quadrantBly float64

	var targetBlx, targetBly float64
	var targetSize float64

	gridCoords = NationalGridSquares[gridRef.Square]

	squareBlx = gridCoords[0] * SquareSize
	squareBly = gridCoords[1] * SquareSize

	switch {
	case gridRef.Quadrant != "":
		subSquareX, _ := strconv.Atoi(string(gridRef.SubSquare[0]))
		subSquareY, _ := strconv.Atoi(string(gridRef.SubSquare[1]))

		subSquareBlx = squareBlx + (float64(subSquareX) * SubSquareSize)
		subSquareBly = squareBly + (float64(subSquareY) * SubSquareSize)

		switch gridRef.Quadrant {
		case SW:
			quadrantBlx = subSquareBlx
			quadrantBly = subSquareBly

		case NW:
			quadrantBlx = subSquareBlx
			quadrantBly = subSquareBly + QuadrantSize

		case SE:
			quadrantBlx = subSquareBlx + QuadrantSize
			quadrantBly = subSquareBly

		case NE:
			quadrantBlx = subSquareBlx + QuadrantSize
			quadrantBly = subSquareBly + QuadrantSize
		}

		targetBlx = quadrantBlx
		targetBly = quadrantBly
		targetSize = QuadrantSize

	case gridRef.SubSquare != "":
		subSquareX, _ := strconv.Atoi(string(gridRef.SubSquare[0]))
		subSquareY, _ := strconv.Atoi(string(gridRef.SubSquare[1]))

		subSquareBlx = squareBlx + (float64(subSquareX) * SubSquareSize)
		subSquareBly = squareBly + (float64(subSquareY) * SubSquareSize)

		targetBlx = subSquareBlx
		targetBly = subSquareBly
		targetSize = SubSquareSize

	case gridRef.Square != "":
		targetBlx = squareBlx
		targetBly = squareBly
		targetSize = SquareSize
	}

	g, err = gridCoordsToGeom(
		[]float64{
			targetBlx,
			targetBly,
		},
		targetSize,
	)
	if err != nil {
		return x, y, err
	}

	center := geom.CenterFromGeometry(g)

	return center[0], center[1], nil
}

func gridCoordsToGeom(bl []float64, tileSize float64) (*geos.Geom, error) {
	var g *geos.Geom

	xmin := bl[0]
	ymin := bl[1]
	xmax := xmin + tileSize
	ymax := ymin + tileSize

	bounds, err := geom.BoundsGeom(xmin, xmax, ymin, ymax)
	if err != nil {
		return g, err
	}

	return bounds, nil
}

func ParseGridRef(ref string) (GridRef, error) {
	var g GridRef

	err := ValidateGridRef(ref)
	if err != nil {
		return g, err
	}

	l := len(ref)

	switch l {
	case 2:
		return GridRef{
			Square: ref,
		}, nil

	case 4:
		square := ref[0:2]
		subsquare := ref[2:4]

		return GridRef{
			Square:    square,
			SubSquare: subsquare,
		}, nil

	case 6:
		square := ref[0:2]
		subsquare := ref[2:4]
		quadrant := ref[4:6]

		return GridRef{
			Square:    square,
			SubSquare: subsquare,
			Quadrant:  Quadrant(quadrant),
		}, nil
	}

	return g, fmt.Errorf("failed to parse %v", ref)
}

func DoOverlap(tl1, br1, tl2, br2 []float64) bool {
	if tl1[0] == br1[0] || tl1[1] == br1[1] || tl2[0] == br2[0] || tl2[1] == br2[1] {
		// the line cannot have positive overlap
		return false
	}

	// if one rectangle is on left side of other
	if tl1[0] >= br2[0] || tl2[0] >= br1[0] {
		return false
	}

	// if one rectangle is above other
	if br1[1] >= tl2[1] || br2[1] >= tl1[1] {
		return false
	}

	return true
}
