package main

import (
	"time"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

func main() {

	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("Voronoi")
	w := 400
	h := 400
	numSeeds := 100
	frameDuration := 0 * time.Millisecond

	v, vErr := NewVoronoi(w, h, numSeeds)
	if vErr != nil {
		panic(vErr)
	}

	g, gErr := NewGame(w, h, v, frameDuration)
	if gErr != nil {
		panic(gErr)
	}

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
