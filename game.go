package main

import (
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type VoronoiDiagram interface {
	Init()
	Tessellate(showIterations bool) error
	ToPixels() []byte
}

// Canvas is the
type Canvas struct {
	width  int
	height int

	gameRunning bool

	voronoi  VoronoiDiagram
	numSeeds int
}

func NewCanvas(
	width int,
	height int,
	voronoi VoronoiDiagram,
) (*Canvas, error) {

	voronoi.Init()

	g := &Canvas{
		width:       width,
		height:      height,
		gameRunning: true,
		voronoi:     voronoi,
	}
	return g, nil
}

func (g *Canvas) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.gameRunning = !g.gameRunning
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.voronoi.Init()
	}

	if g.gameRunning {
		return g.voronoi.Tessellate(true)
	}
	return nil
}

func (g *Canvas) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.voronoi.ToPixels())
}

func (g *Canvas) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}
