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

func Generate(seed int64) *[64][64]protocol.Surface_Type {
	rand.Seed(seed)
	var gradient [64][64]int64
	step := 32 // partition size
	iter := int64(2)
	// Fill the gradient in increasingly smaller grains.
	for step > 1 {
		delta := math.MaxInt64 / iter
		for x := 0; x < 64; x += step {
			for y := 0; y < 64; y += step {
				// Randomly choose up or down.
				polarity := rand.Intn(2)
				if polarity == 0 {
					delta = -delta
				}
				for i := x; i < x+step-1; i++ {
					for j := y; j < y+step-1; j++ {
						// Apply the delta amortized by
						// proximity to partition center.
						relativeDelta := int64(float64(delta) * (float64(i-step/2) / float64(step)))
						gradient[i][j] += relativeDelta
					}
				}
			}
		}
		// As the partitions become smaller,
		// so do the vertical deltas.
		step = step / 2
		iter = iter * 2
	}
	// Map integer gradients to surface types.
	var surface [64][64]protocol.Surface_Type
	for x := 0; x < 64; x++ {
		for y := 0; y < 64; y++ {
			i := gradient[x][y]
			switch {
			case i < math.MaxInt64/-2:
				surface[x][y] = protocol.Surface_WATER
			case i < 0:
				surface[x][y] = protocol.Surface_GRASS
			case i < math.MaxInt64/2:
				surface[x][y] = protocol.Surface_DIRT
			default:
				surface[x][y] = protocol.Surface_ROCK
			}
		}
	}
	return &surface
}
