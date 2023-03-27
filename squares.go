package nationalgrid

// the dimension of the source (osdata) shapefile (100KM).
const (
	osShpFileSize = float64(100000)
	SquareSize    = osShpFileSize
	SubSquareSize = SquareSize / 10
	QuadrantSize  = SubSquareSize / 2
)

var NationalGridSquares = map[string][]float64{
	"HP": {
		4,
		12,
	},
	"HT": {
		3,
		11,
	},
	"HU": {
		4,
		11,
	},
	"HW": {
		1,
		10,
	},
	"HX": {
		2,
		10,
	},
	"HY": {
		3,
		10,
	},
	"HZ": {
		4,
		10,
	},
	"NA": {
		0,
		9,
	},
	"NB": {
		1,
		9,
	},
	"NC": {
		2,
		9,
	},
	"ND": {
		3,
		9,
	},
	"NF": {
		0,
		8,
	},
	"NG": {
		1,
		8,
	},
	"NH": {
		2,
		8,
	},
	"NJ": {
		3,
		8,
	},
	"NK": {
		4,
		8,
	},
	"NL": {
		0,
		7,
	},
	"NM": {
		1,
		7,
	},
	"NN": {
		2,
		7,
	},
	"NO": {
		3,
		7,
	},
	"NR": {
		1,
		6,
	},
	"NS": {
		2,
		6,
	},
	"NT": {
		3,
		6,
	},
	"NU": {
		4,
		6,
	},
	"NW": {
		1,
		5,
	},
	"NX": {
		2,
		5,
	},
	"NY": {
		3,
		5,
	},
	"NZ": {
		4,
		5,
	},
	"OV": {
		5,
		5,
	},
	"SD": {
		3,
		4,
	},
	"SE": {
		4,
		4,
	},
	"TA": {
		5,
		4,
	},
	"SH": {
		2,
		3,
	},
	"SJ": {
		3,
		3,
	},
	"SK": {
		4,
		3,
	},
	"TF": {
		5,
		3,
	},
	"TG": {
		6,
		3,
	},
	"SM": {
		1,
		2,
	},
	"SN": {
		2,
		2,
	},
	"SO": {
		3,
		2,
	},
	"SP": {
		4,
		2,
	},
	"TL": {
		5,
		2,
	},
	"TM": {
		6,
		2,
	},
	"SR": {
		1,
		1,
	},
	"SS": {
		2,
		1,
	},
	"ST": {
		3,
		1,
	},
	"SU": {
		4,
		1,
	},
	"TQ": {
		5,
		1,
	},
	"TR": {
		6,
		1,
	},
	"SV": {
		0,
		0,
	},
	"SW": {
		1,
		0,
	},
	"SX": {
		2,
		0,
	},
	"SY": {
		3,
		0,
	},
	"SZ": {
		4,
		0,
	},
	"TV": {
		5,
		0,
	},
}
