package nationalgrid

import (
	"fmt"

	"github.com/wroge/wgs84"
)

type LocationType int

const (
	WGS84 LocationType = iota
	OSGB36
	NATIONALGRID
)

func (l LocationType) String() string {
	return [...]string{"WGS84", "OSGB36", "NATIONALGRID"}[l]
}

// Longitude is X (Easting), Latitude is Y (northing)
// GridRef is eg. SD, or SD91SW

// WGS84 - Lon, Lat
// OSGB36 - X, Y
// NATIONALGRID - GridRef

type LonLat struct {
	Lon float64
	Lat float64
}

type EastNorth struct {
	X float64 //East
	Y float64 //North
}

type Location struct {
	Type      string
	LonLat    LonLat
	EastNorth EastNorth
	GridRef   string
}

func (c LonLat) ToString() string {
	return fmt.Sprintf("%+v", c)
}

func (c EastNorth) ToString() string {
	return fmt.Sprintf("%+v", c)
}

func (c Location) ToOSGB36() EastNorth {
	var r EastNorth
	h := 0.0

	switch c.Type {
	case WGS84.String():
		east, north, _ := wgs84.LonLat().To(wgs84.OSGB36NationalGrid())(c.LonLat.Lon, c.LonLat.Lat, h)
		r = EastNorth{
			X: east,
			Y: north,
		}
	case OSGB36.String():
		r = c.EastNorth
	case NATIONALGRID.String():
		east, north, _ := GetGridEastNorth(c.GridRef)
		r = EastNorth{
			X: north,
			Y: east,
		}
	}

	return r
}
