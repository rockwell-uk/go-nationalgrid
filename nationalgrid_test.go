package nationalgrid

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"path"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/rockwell-uk/go-geos-draw/geom"
	"github.com/rockwell-uk/go-text/fonts"
	geos "github.com/twpayne/go-geos"
)

var (
	gctx     = geos.NewContext()
	fontData = draw2d.FontData{
		Name:   "bold",
		Family: draw2d.FontFamilySans,
		Style:  draw2d.FontStyleNormal,
	}
	textRotation = 0.0
	black        = color.RGBA{0x00, 0x00, 0x00, 0xFF}
	white        = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	scale        = func(x, y float64) (float64, float64) {
		return x, y
	}
	fillColor   = white
	strokeWidth = 0.0
	strokeColor = black
	lineWidth   = 1.0
)

func TestLogSquareCentres(t *testing.T) {
	outFileName := "test-output/national-grid-squares.txt"
	os.Remove(outFileName)

	f, err := os.OpenFile(outFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// sort alphabetically
	keys := make([]string, 0, len(NationalGridSquares))
	for ref := range NationalGridSquares {
		keys = append(keys, ref)
	}
	sort.Strings(keys)

	for _, ref := range keys {
		gridRef, err := ParseGridRef(ref)
		if err != nil {
			t.Fatal(err)
		}

		east, north, err := getGridCoordCenter(gridRef)
		if err != nil {
			t.Fatal(err)
		}

		n := strconv.FormatFloat(north, 'f', -1, 64)
		e := strconv.FormatFloat(east, 'f', -1, 64)

		line := fmt.Sprintf("%v [%v, %v]", ref, e, n)

		if _, err := f.WriteString(line + "\n"); err != nil {
			t.Fatal(err)
		}
	}
}

func TestLogSubSquareCentres(t *testing.T) {
	outFileName := "test-output/national-grid-subsquares.txt"
	os.Remove(outFileName)

	f, err := os.OpenFile(outFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// sort alphabetically
	keys := make([]string, 0, len(NationalGridSquares))
	for ref := range NationalGridSquares {
		keys = append(keys, ref)
	}
	sort.Strings(keys)

	for _, ref := range keys {
		for i := 0; i <= 99; i++ {
			subsquare := fmt.Sprintf("%02d", i)
			ssRef := fmt.Sprintf("%v%v", ref, subsquare)

			gridRef, err := ParseGridRef(ssRef)
			if err != nil {
				t.Fatal(err)
			}

			east, north, err := getGridCoordCenter(gridRef)
			if err != nil {
				t.Fatal(err)
			}

			n := strconv.FormatFloat(north, 'f', -1, 64)
			e := strconv.FormatFloat(east, 'f', -1, 64)

			line := fmt.Sprintf("%v [%v, %v]", ssRef, e, n)

			if _, err := f.WriteString(line + "\n"); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestParseGridRef(t *testing.T) {
	tests := map[string]struct {
		Ref  GridRef
		Fail bool
	}{
		"SD": {
			Ref: GridRef{
				Square: "SD",
			},
		},
		"SD00": {
			Ref: GridRef{
				Square:    "SD",
				SubSquare: "00",
			},
		},
		"SD00SE": {
			Ref: GridRef{
				Square:    "SD",
				SubSquare: "00",
				Quadrant:  Quadrant("SE"),
			},
		},
		"SD00NW": {
			Ref: GridRef{
				Square:    "SD",
				SubSquare: "00",
				Quadrant:  Quadrant("NW"),
			},
		},
		"SD00XX": {
			Fail: true,
		},
		"00XX": {
			Fail: true,
		},
	}

	for ref, tt := range tests {
		actual, err := ParseGridRef(ref)

		if !tt.Fail && err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(tt.Ref, actual) {
			t.Fatalf("%v\nexpected %#v\ngot %#v", ref, tt.Ref, actual)
		}
	}
}

func TestGetGridLatLonSubSquare(t *testing.T) {
	gridRef := "SD00"

	expectedeast := 305000.0
	expectednorth := 405000.0

	east, north, err := GetGridLatLon(gridRef)
	if err != nil {
		t.Fatal(err)
	}

	if expectedeast != east {
		t.Fatalf("expectedeast %+v, got %+v", expectedeast, east)
	}

	if expectednorth != north {
		t.Fatalf("expectednorth %+v, got %+v", expectednorth, north)
	}
}

func TestGetGridLatLonQuadrant(t *testing.T) {
	tests := map[string]struct {
		Expected []float64
	}{
		"SV00SW": {
			Expected: []float64{
				2500.0,
				2500.0,
			},
		},
		"SD00SW": {
			Expected: []float64{
				302500.0,
				402500.0,
			},
		},
		"SD91NW": {
			Expected: []float64{
				392500.0,
				417500.0,
			},
		},
		"SV00": {
			Expected: []float64{
				5000.0,
				5000.0,
			},
		},
		"SD01": {
			Expected: []float64{
				305000.0,
				415000.0,
			},
		},
		"SD91": {
			Expected: []float64{
				395000.0,
				415000.0,
			},
		},
		"SV": {
			Expected: []float64{
				50000.0,
				50000.0,
			},
		},
		"SD": {
			Expected: []float64{
				350000.0,
				450000.0,
			},
		},
	}

	for ref, tt := range tests {
		east, north, err := GetGridLatLon(ref)
		if err != nil {
			t.Fatal(err)
		}

		actual := []float64{
			east,
			north,
		}
		if !reflect.DeepEqual(tt.Expected, actual) {
			t.Fatalf("%v expected %+v, got %+v", ref, tt.Expected, actual)
		}
	}
}

func TestGetGridLatLonSquare(t *testing.T) {
	gridRef := "SD"
	err := ValidateSquare(gridRef)
	if err != nil {
		t.Fatal(err)
	}

	expectedeast := 350000.0
	expectednorth := 450000.0

	east, north, err := GetGridLatLon(gridRef)
	if err != nil {
		t.Fatal(err)
	}

	if expectedeast != east {
		t.Fatalf("expectedeast %+v, got %+v", expectedeast, east)
	}

	if expectednorth != north {
		t.Fatalf("expectednorth %+v, got %+v", expectednorth, north)
	}
}

func TestGetGridCoordCenter(t *testing.T) {
	ref := "SD"
	expectedeast := 350000.0
	expectednorth := 450000.0

	gridRef, err := ParseGridRef(ref)
	if err != nil {
		t.Fatal(err)
	}
	east, north, err := getGridCoordCenter(gridRef)
	if err != nil {
		t.Fatal(err)
	}

	if expectedeast != east {
		t.Fatalf("expectedeast %+v, got %+v", expectedeast, east)
	}

	if expectednorth != north {
		t.Fatalf("expectednorth %+v, got %+v", expectednorth, north)
	}
}

func TestGridCoordsToGeom(t *testing.T) {
	sdwkt := "POLYGON ((300000.0000000000000000 400000.0000000000000000, 400000.0000000000000000 400000.0000000000000000, 400000.0000000000000000 500000.0000000000000000, 300000.0000000000000000 500000.0000000000000000, 300000.0000000000000000 400000.0000000000000000))"

	expected, err := gctx.NewGeomFromWKT(sdwkt)
	if err != nil {
		t.Fatal(err)
	}

	gridRef := "SD"
	err = ValidateSquare(gridRef)
	if err != nil {
		t.Fatal(err)
	}

	gridCoords := NationalGridSquares[gridRef]

	g, err := gridCoordsToGeom(gridCoords, SquareSize)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, g) {
		t.Fatalf("expected %+v, got %+v", expected, g)
	}
}

func TestGetSubSquares(t *testing.T) {
	expected := map[string][]int{
		"sd": {
			81,
			91,
		},
	}

	targetTilePoly := "POLYGON ((387221.1985319799860008 410715.0784210899728350, 392221.1985319799860008 410715.0784210899728350, 392221.1985319799860008 415715.0784210899728350, 387221.1985319799860008 415715.0784210899728350, 387221.1985319799860008 410715.0784210899728350))"
	targetTileGeom, _ := gctx.NewGeomFromWKT(targetTilePoly)

	actual := GetSubSquares(targetTileGeom.Bounds())

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %+v, got %+v", expected, actual)
	}
}

func TestGetAllFullSquareCoords(t *testing.T) {
	expected := []int{}
	for i := 0; i <= 99; i++ {
		expected = append(expected, i)
	}

	tests := map[string]*geos.Bounds{}

	for key, square := range NationalGridSquares {
		tileX := square[0] * SquareSize
		tileY := square[1] * SquareSize

		minX := tileX
		minY := tileY
		maxX := tileX + SquareSize
		maxY := tileY + SquareSize

		env := geos.NewBounds(minX, minY, maxX, maxY)

		k := strings.ToLower(key)
		tests[k] = env
	}

	for tile, envelope := range tests {
		actual := GetSubSquares(envelope)

		if !reflect.DeepEqual(expected, actual[tile]) {
			t.Fatalf("\nexpected %+v\ngot %+v", expected, actual[tile])
		}
	}
}

func TestDrawSquares(t *testing.T) {
	imgWidth := 700
	imgHeight := 1300

	m, gc := setupImage(imgWidth, imgHeight)

	fontSize := 10.0
	typeFace := getTypeFace(gc, fontSize)
	fonts.SetFont(gc, typeFace)

	tileSize := SquareSize / 1000

	for label, square := range NationalGridSquares {
		tileX := square[0] * tileSize
		tileY := square[1] * tileSize

		xPos := tileX
		yPos := float64(imgHeight) - tileY - tileSize

		g, err := geom.BoundsGeom(
			xPos,
			xPos+tileSize,
			yPos,
			yPos+tileSize,
		)
		if err != nil {
			t.Fatal(err)
		}

		l, err := geom.ToLineString(g)
		if err != nil {
			t.Fatal(err)
		}

		err = geom.DrawLine(gc, l, lineWidth, fillColor, strokeWidth, strokeColor, scale)
		if err != nil {
			t.Fatal(err)
		}

		textWidth := fonts.GetTextWidth(typeFace, label)

		labelPos := []float64{
			xPos + (SquareSize * 5 / 10000) - textWidth/2,
			yPos + (SquareSize * 5 / 10000),
		}

		err = geom.DrawString(gc, labelPos, textRotation, label)
		if err != nil {
			t.Fatal(err)
		}
	}

	// draw the image
	err := savePNG("test-output/all.png", m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDrawSubSectors(t *testing.T) {
	imgWidth := 1000
	imgHeight := 1000

	sectorName := "SV"
	if _, ok := NationalGridSquares[sectorName]; !ok {
		t.Fatalf("Unable to load sector %v", sectorName)
	}

	m, gc := setupImage(imgWidth, imgHeight)

	fontSize := 30.0
	typeFace := getTypeFace(gc, fontSize)
	fonts.SetFont(gc, typeFace)
	fm := fonts.GetFaceMetrics(typeFace)
	textWidth := fonts.GetTextWidth(typeFace, sectorName)

	tileSize := SquareSize / 100

	if _, ok := NationalGridSquares[sectorName]; !ok {
		t.Fatalf("Unable to load sector %v", sectorName)
	}
	sector := NationalGridSquares[sectorName]

	tileX := sector[0] * SquareSize
	tileY := sector[1] * SquareSize

	tileX /= 1000
	tileY /= 1000

	xPos := tileX
	yPos := float64(imgHeight) - tileY - tileSize

	g, err := geom.BoundsGeom(
		xPos,
		xPos+tileSize,
		yPos,
		yPos+tileSize,
	)
	if err != nil {
		t.Fatal(err)
	}

	l, err := geom.ToLineString(g)
	if err != nil {
		t.Fatal(err)
	}

	err = geom.DrawLine(gc, l, lineWidth, fillColor, strokeWidth, strokeColor, scale)
	if err != nil {
		t.Fatal(err)
	}

	labelPos := []float64{
		xPos + (SquareSize * 5 / 1000) - textWidth/2,
		yPos + (SquareSize * 5 / 1000) + ((fm.Ascent - fm.Descent) / 2),
	}

	err = geom.DrawString(gc, labelPos, textRotation, sectorName)
	if err != nil {
		t.Fatal(err)
	}

	// subsectors
	fontSize = 10.0
	typeFace = getTypeFace(gc, fontSize)
	fonts.SetFont(gc, typeFace)
	fm = fonts.GetFaceMetrics(typeFace)

	for i := 0; i <= 99; i++ {
		subsectorSize := tileSize / 10

		bits := fmt.Sprintf("%02d", i)
		addX, _ := strconv.Atoi(string(bits[0]))
		addY, _ := strconv.Atoi(string(bits[1]))

		blX := tileX + (float64(addX) * subsectorSize)
		blY := float64(imgHeight) - subsectorSize - (tileY + (float64(addY) * subsectorSize))

		textWidth := fonts.GetTextWidth(typeFace, bits)

		g, err := geom.BoundsGeom(
			blX,
			blX+subsectorSize,
			blY,
			blY+subsectorSize,
		)
		if err != nil {
			t.Fatal(err)
		}

		l, err := geom.ToLineString(g)
		if err != nil {
			t.Fatal(err)
		}

		err = geom.DrawLine(gc, l, lineWidth, fillColor, strokeWidth, strokeColor, scale)
		if err != nil {
			t.Fatal(err)
		}

		labelPos := []float64{
			blX + (subsectorSize / 2) - (textWidth / 2),
			blY + (subsectorSize / 2) + ((fm.Ascent - fm.Descent) / 2),
		}

		err = geom.DrawString(gc, labelPos, textRotation, bits)
		if err != nil {
			t.Fatal(err)
		}
	}

	// draw the image
	err = savePNG("test-output/subsectors.png", m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDoOverlap(t *testing.T) {
	tests := []struct {
		tl1      []float64
		br1      []float64
		tl2      []float64
		br2      []float64
		overlaps bool
	}{
		{
			[]float64{0, 2},
			[]float64{1, 1},
			[]float64{1, 1},
			[]float64{0, 2},
			false,
		},
		{
			[]float64{0, 3},
			[]float64{2, 1},
			[]float64{1, 2},
			[]float64{3, 0},
			true,
		},
	}

	for _, tt := range tests {
		doOverlap := DoOverlap(tt.tl1, tt.br1, tt.tl2, tt.br2)

		if doOverlap != tt.overlaps {
			t.Fatalf("expected %+v, got %+v", tt.overlaps, doOverlap)
		}
	}
}

func setupImage(width, height int) (*image.RGBA, *draw2dimg.GraphicContext) {
	m := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(m, m.Bounds(), &image.Uniform{white}, image.Point{0, 0}, draw.Src)
	gc := draw2dimg.NewGraphicContext(m)

	gc.SetDPI(72)

	return m, gc
}

func getTypeFace(gc *draw2dimg.GraphicContext, fontSize float64) fonts.TypeFace {
	strokeStyle := draw2d.StrokeStyle{
		Color: white,
		Width: lineWidth,
	}
	return fonts.TypeFace{
		StrokeStyle: strokeStyle,
		Color:       black,
		Size:        fontSize,
		FontData:    fontData,
		Face:        fonts.GetFace(gc, fontData, fontSize),
	}
}

func savePNG(fname string, m image.Image) error {
	dir, _ := path.Split(fname)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	err = draw2dimg.SaveToPngFile(fname, m)
	if err != nil {
		return err
	}

	return nil
}
