package main

import (
	"errors"
	"math/rand"
	"time"
)

// Voronoi is the engine used to generate a voronoi diagram on a canvas, starting from auto-generated seed points
type Voronoi struct {

	// diagram size (in pixels)
	width  int
	height int

	// seed configuration of the diagram
	numSeeds int     // number of seeds for the diagram
	seeds    []Point // list of seeds for the diagram

	radius      int     // current radius of the computation
	activeSeeds []Point // list of active seeds to take into account for the computation

	distances [][]int // precomputed distances matrix (for efficiency reasons)

	diagram [][]*Point // resulting diagram (initially empty, to be computed)
}

// NewVoronoi creates a new diagram struct
func NewVoronoi(
	width int,
	height int,
	numSeeds int,
) (*Voronoi, error) {

	if numSeeds > width*height {
		return nil, errors.New("Number of seeds cannot be more than the pixels in the canvas")
	}

	return &Voronoi{
		width:       width,
		height:      height,
		numSeeds:    numSeeds,
		seeds:       []Point{},
		radius:      0,
		activeSeeds: []Point{},
		distances:   make([][]int, 2*width+1),
		diagram:     make([][]*Point, width),
	}, nil
}

// Init initializes the Voronoi diagram and generates a new set of seeds
func (v *Voronoi) Init() {
	v.initDistances()
	v.initDiagram()
	v.initSeeds()
	v.initTessellation()
}

// initDistances populates the precomputed distances matrix,
// to avoid recomputing the same distance values over and over
func (v *Voronoi) initDistances() {

	// the distance vectors needed by the engine can assume values up to twice their dimension  (2*width or 2*height)
	for i := 0; i <= 2*v.width; i++ {

		column := make([]int, 2*v.height+1)
		v.distances[i] = column

		for j := 0; j <= 2*v.height; j++ {
			v.distances[i][j] = i*i + j*j
		}
	}
}

// initDiagram populates the diagram with empty points
func (v *Voronoi) initDiagram() {

	for i := 0; i < v.width; i++ {

		column := make([]*Point, v.height)
		v.diagram[i] = column

		for j := 0; j < v.height; j++ {
			v.diagram[i][j] = nil
		}
	}
}

// initSeeds generates a random set of seeds with random colors and stores them in the diagram
func (v *Voronoi) initSeeds() {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	v.seeds = []Point{}

	for i := 0; i < v.numSeeds; i++ {
		x := int(r.Intn(v.width))
		y := int(r.Intn(v.height))
		d := 0
		seed := Point{
			X:        x,
			Y:        y,
			Distance: &d,
			Color: &Color{
				R: uint8(r.Intn(256)),
				G: uint8(r.Intn(256)),
				B: uint8(r.Intn(256)),
				A: uint8(r.Intn(256)),
			},
		}

		v.seeds = append(v.seeds, seed)
		v.diagram[seed.X][seed.Y] = &seed
	}
}

// initTessellation starts the tessellation of the existing set of seeds
func (v *Voronoi) initTessellation() {

	v.radius = 0
	v.activeSeeds = v.seeds

	// fmt.Println("#######################################")
	// fmt.Println("#### Voronoi tessellation starting ####")
	// fmt.Println("#######################################")
}

/*
	Tessellate computes the voronoi diagram

	It works on a list of 'active' seeds, where 'active' means that the seed can still extend its area.
	At each iteration,
*/
func (v *Voronoi) Tessellate(hideIterations bool) error {

	for len(v.activeSeeds) > 0 {

		stillActiveSeeds := []Point{}
		incrementalVectors := v.getIncrementalVectors()

		for _, seed := range v.activeSeeds {
			// fmt.Println("Iteration starting. Active seeds: ", len(v.activeSeeds))

			stillActive := false
			for _, incrementalVector := range incrementalVectors {
				stillActive = v.assignPointToSeed(
					seed,
					v.distances[abs(incrementalVector.X)][abs(incrementalVector.Y)],
					incrementalVector.X,
					incrementalVector.Y,
				) || stillActive
			}

			if stillActive {
				stillActiveSeeds = append(stillActiveSeeds, seed)
			}
		}

		v.activeSeeds = stillActiveSeeds

		if !hideIterations {
			break
		}
	}

	return nil
}

func (v *Voronoi) assignPointToSeed(seed Point, distance int, dx int, dy int) bool {
	if seed.X+dx < 0 ||
		seed.X+dx >= v.width ||
		seed.Y+dy < 0 ||
		seed.Y+dy >= v.height {
		// fmt.Println(fmt.Sprintf("Point (%d,%d) out of canvas, discarded", seed.X+dx, seed.Y+dy))
		return false
	}

	p := v.pointFromDiagram(seed.X+dx, seed.Y+dy)

	if p.Distance != nil && *p.Distance < distance {
		// fmt.Println(fmt.Sprintf("Point (%d,%d) has already a smaller distance (%d < %d), discarded", seed.X+dx, seed.Y+dy, *p.Distance, distance))
		return false
	}

	// fmt.Println(fmt.Sprintf("Assigning point (%d,%d) to cell with seed (%d, %d). Distance: %d", p.X, p.Y, seed.X, seed.Y, distance))
	p.Color = seed.Color
	p.Distance = &distance
	v.diagram[p.X][p.Y] = &p

	return true
}

/*
	getIncrementalVectors

	It returns a list of points, intended as coordinates relative to the seed,
	that represents the new layer of pixels of the expanding cell.
*/
func (v *Voronoi) getIncrementalVectors() []Point {
	combinations := []Point{}

	v.radius++ // increment the radius of the cell

	dx := 0
	dy := v.radius

	for dy >= dx {
		combinations = append(combinations, Point{X: dx, Y: dy})
		combinations = append(combinations, Point{X: dx, Y: -dy})
		combinations = append(combinations, Point{X: -dx, Y: dy})
		combinations = append(combinations, Point{X: -dx, Y: -dy})
		combinations = append(combinations, Point{X: dy, Y: dx})
		combinations = append(combinations, Point{X: dy, Y: -dx})
		combinations = append(combinations, Point{X: -dy, Y: dx})
		combinations = append(combinations, Point{X: -dy, Y: -dx})

		dx++
		dy--
	}
	return combinations
}

// pointFromDiagram gets the point of the diagram corresponding to the given coordinates
func (v *Voronoi) pointFromDiagram(x int, y int) Point {
	if v.diagram[x][y] == nil {
		v.diagram[x][y] = &Point{
			X: x,
			Y: y,
		}
	}

	return *v.diagram[x][y]
}

// ToPixels generates the byte array containing the information to render the diagram.
// Each row of the canvas is concatenated to obtain a one-dimensional array.
// Each pixel is represented by 4 bytes, representing the Red, Green, Blue and Alpha info.
func (v *Voronoi) ToPixels() []byte {
	pixels := make([]byte, v.width*v.height*4)

	// iterate through each pixel
	for i := 0; i < v.width; i++ {
		for j := 0; j < v.height; j++ {
			pos := (j*v.width + i) * 4

			if v.diagram[i][j] != nil && v.diagram[i][j].Color != nil {
				pixels[pos] = v.diagram[i][j].Color.R
				pixels[pos+1] = v.diagram[i][j].Color.G
				pixels[pos+2] = v.diagram[i][j].Color.B
				pixels[pos+3] = v.diagram[i][j].Color.A

			} else {
				// if the point has not assigned any color yet, show it as black
				pixels[pos] = 0
				pixels[pos+1] = 0
				pixels[pos+2] = 0
				pixels[pos+3] = 0
			}
		}
	}

	// iterate through the seeds to render them as black points
	for _, s := range v.seeds {
		pos := (s.Y*v.width + s.X) * 4
		pixels[pos] = 0
		pixels[pos+1] = 0
		pixels[pos+2] = 0
		pixels[pos+3] = 0
	}

	return pixels
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
