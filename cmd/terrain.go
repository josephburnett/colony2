package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/josephburnett/colony2/pkg/protocol"
)

func main() {
	v := &protocol.View{
		World: &protocol.World{
			Surfaces: make(map[int32]*protocol.World_SurfaceRow),
		},
		XMin: 0,
		YMin: 0,
		XMax: 64,
		YMax: 64,
	}
	surface := Generate(int64(999))
	for x := 0; x < 64; x++ {
		row, ok := v.World.Surfaces[int32(x)]
		if !ok {
			row = &protocol.World_SurfaceRow{
				Columns: make(map[int32]*protocol.Surface),
			}
			v.World.Surfaces[int32(x)] = row
		}
		for y := 0; y < 64; y++ {
			row.Columns[int32(y)] = &protocol.Surface{
				Type: surface[x][y],
			}
		}
	}
	fmt.Printf("%v", v.Render())
}

func Generate(xMin, yMin, xMax, yMax int32) *[64][64]protocol.Surface_Type {

	var gradient [64][64]float64
	var r *rand.Rand

	// Fill the corners with initial values.
	for _, c := range [][2]int32{
		{xMax, yMax}, // upper right
		{xMin, yMax}, // upper left
		{xMax, yMin}, // lower right
		{xMin, yMin}, // lower left
	} {
		x := int64(c[0])
		y := int64(c[1])

		// Seed with coordinates to get a deterministic height.
		seed := x<<32 + y
		r = rand.New(rand.NewSource(seed))
		height := r.Float64()
		gradient[int(x)][int(y)] = height
	}

	// Run the diamond-square algorithm to fill in the rest of the points.
	// See https://en.wikipedia.org/wiki/Diamond-square_algorithm.
	fill := func(x1, y1, x2, y2 int32, depth int) {

		// Diamond step.
		avg := (gradient[x1][y1] + gradient[x1][y2] +
			gradient[x2][y] + gradient[x2][y2]) / 4.0
		// Random delta, getting exponentially smaller.
		delta := (r.Float64() - 0.5) * (1.0 / math.Exp2(float64(depth)))
		centerX := (x1 + x2) / 2
		centerY := (y1 + y2) / 2
		gradient[centerX][centerY] = avg + delta

		// Square step.
		wrapping := func(x, y int32) float64 {
			// TODO: wrap proportially with depth
			if x < xMin {
				x = xMax
			}
			if x > xMax {
				x = xMin
			}
			if y < yMin {
				y = yMax
			}
			if y > yMax {
				y = yMin
			}
			return gradient[x][y]
		}
		// TODO: Diamond step.
	}
}
