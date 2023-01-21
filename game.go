package main

import (
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// VoronoiDiagram is the voronoi engine
type VoronoiDiagram interface {
	Init()
	Tessellate(hideIterations bool) error
	ToPixels() []byte
}

// Canvas handles the canvas visualization
type Canvas struct {

	// resolution of the canvas
	width  int
	height int

	gameRunning    bool
	hideIterations bool

	voronoi VoronoiDiagram
}

// NewCanvas creates a canvas with a voronoi ready to start
func NewCanvas(
	width int,
	height int,
	hideIterations bool,
	voronoi VoronoiDiagram,
) (*Canvas, error) {

	voronoi.Init()

	g := &Canvas{
		width:          width,
		height:         height,
		gameRunning:    true,
		hideIterations: hideIterations,
		voronoi:        voronoi,
	}
	return g, nil
}

// Update computes a new frame
func (g *Canvas) Update() error {

	// Intercepts the Enter key and starts/stops the execution
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.gameRunning = !g.gameRunning
	}

	// Intercepts the Space key
	// and restarts the execution regenerating the seeds
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.voronoi.Init()
	}

	if g.gameRunning {
		// compute the voronoi tessellation
		return g.voronoi.Tessellate(g.hideIterations)
	}
	return nil
}

// Draw writes the computed frame as a byte sequence
func (g *Canvas) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.voronoi.ToPixels())
}

// Layout returns the resolution of the canvas
func (g *Canvas) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}
