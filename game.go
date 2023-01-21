package main

import (
	"time"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type VoronoiDiagram interface {
	Init()
	Tessellate(showIterations bool) error
	ToPixels() []byte
}

type Game struct {
	width  int
	height int

	frameDuration time.Duration
	gameRunning   bool

	numSeeds int
	voronoi  VoronoiDiagram
}

func NewGame(
	width int,
	height int,
	voronoi VoronoiDiagram,
	frameDuration time.Duration,
) (*Game, error) {

	voronoi.Init()

	g := &Game{
		width:         width,
		height:        height,
		frameDuration: frameDuration,
		gameRunning:   true,
		voronoi:       voronoi,
	}
	return g, nil
}

func (g *Game) Update() error {
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

func (g *Game) Draw(screen *ebiten.Image) {
	if g.frameDuration > 0 {
		time.Sleep(g.frameDuration)
	}
	screen.WritePixels(g.voronoi.ToPixels())
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}
