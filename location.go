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

type LatLon struct {
	Lat float64
	Lon float64
}

type Location struct {
	Type    string
	LatLon  LatLon
	GridRef string
}

type OSGB36LatLon struct {
	LatLon LatLon
}

func (c OSGB36LatLon) ToString() string {
	return fmt.Sprintf("%+v", c)
}

func (c Location) ToOSGB36() OSGB36LatLon {
	var r OSGB36LatLon
	h := 0.0

	switch c.Type {
	case WGS84.String():
		east, north, _ := wgs84.LonLat().To(wgs84.OSGB36NationalGrid())(c.LatLon.Lon, c.LatLon.Lat, h)
		r = OSGB36LatLon{
			LatLon: LatLon{
				Lat: east,
				Lon: north,
			},
		}
	case OSGB36.String():
		r = OSGB36LatLon{
			LatLon: LatLon{
				Lat: c.LatLon.Lat,
				Lon: c.LatLon.Lon,
			},
		}
	case NATIONALGRID.String():
		east, north, _ := GetGridLatLon(c.GridRef)
		r = OSGB36LatLon{
			LatLon: LatLon{
				Lat: east,
				Lon: north,
			},
		}
	}

	return r
}
