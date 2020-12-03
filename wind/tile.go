package wind

import (
	"image"
	"image/color"
	"log"
	"math"
	"time"
)

var colors map[float64]([3]float64) = map[float64]([3]float64){
	0:   [3]float64{98, 113, 184},
	2.5: [3]float64{61, 110, 163},
	5:   [3]float64{74, 148, 170},
	7.5: [3]float64{74, 146, 148},
	10:  [3]float64{77, 142, 124},
	15:  [3]float64{76, 164, 76},
	20:  [3]float64{103, 164, 54},
	25:  [3]float64{162, 135, 64},
	30:  [3]float64{162, 109, 92},
	35:  [3]float64{141, 63, 92},
	40:  [3]float64{151, 75, 145},
	50:  [3]float64{95, 100, 160},
	60:  [3]float64{91, 136, 161},
}

var keys []float64 = []float64{0, 2.5, 5, 7.5, 10, 15, 20, 25, 30, 35, 40, 50, 60}

func GenerateTile(w *Winds, z, y, x int, m time.Time) *image.RGBA {
	w1, w2, h := w.FindWinds(m)

	tile := image.NewRGBA(image.Rect(0, 0, 256, 256))

	bb := tile2boudingbox(z, x, y)

	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {

			lat := bb.north + float64(j)*(bb.south-bb.north)/256
			lon := bb.west + float64(i)*(bb.east-bb.west)/256

			_, d := Interpolate(w1, w2, lat, lon, h)
			d *= 1.9438444924406

			k := 0
			key := keys[k]
			for k, key = range keys {
				if key < d {
					continue
				} else {
					break
				}
			}

			kh := 0.0
			key1 := keys[k]
			key2 := keys[k]
			if key2 > d && k > 0 {
				key1 = keys[k-1]
				kh = (d - key1) / (key2 - key1)
			}

			r := uint8(colors[key1][0]*(1-kh) + colors[key2][0]*kh)
			g := uint8(colors[key1][1]*(1-kh) + colors[key2][1]*kh)
			b := uint8(colors[key1][2]*(1-kh) + colors[key2][2]*kh)

			// log.Println(d, key1, key2, kh, r, g, b)

			tile.Set(i, j, color.RGBA{r, g, b, 0})
		}
	}

	if w1 != nil {
		log.Println("w1", w1[0].File)
	}
	if w2 != nil {
		log.Println("w2", w2[0].File)
	}
	log.Println("h", h)

	return tile
}

type boundingbox struct {
	north float64
	south float64
	east  float64
	west  float64
}

func tile2boudingbox(z, x, y int) boundingbox {
	return boundingbox{
		north: tile2lat(z, y),
		south: tile2lat(z, y+1),
		west:  tile2lon(z, x),
		east:  tile2lon(z, x+1),
	}
}

func tile2lon(z, x int) float64 {
	return float64(x)/math.Exp2(float64(z))*360.0 - 180.0
}

func tile2lat(z, y int) float64 {
	n := math.Pi - 2.0*math.Pi*float64(y)/math.Exp2(float64(z))
	return 180.0 / math.Pi * math.Atan(math.Sinh(n))
}
