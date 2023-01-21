package main

import (
	"errors"
	"math/rand"
	"time"
)

// Voronoi is used to generate a voronoi diagram on a canvas, starting from auto-generated seed points
type Voronoi struct {
	// input parameters
	width    int
	height   int
	numSeeds int

	// list of seeds for the Voronoi diagram
	seeds       []Point
	activeSeeds []Point

	// resulting diagram (to be constructed)
	diagram [][]*Point

	radius    int
	distances [][]int
}

// NewVoronoi creates a new diagram struct
func NewVoronoi(width int, height int, numSeeds int) (*Voronoi, error) {

	if numSeeds > width*height {
		return nil, errors.New("Number of seeds cannot be more than the pixels in the canvas")
	}

	return &Voronoi{
		width:       width,
		height:      height,
		distances:   make([][]int, 2*width+1),
		numSeeds:    numSeeds,
		seeds:       []Point{},
		activeSeeds: []Point{},
		diagram:     make([][]*Point, width),
	}, nil
}

// Init initializes the Voronoi diagram and generates a set of seeds
func (v *Voronoi) Init() {
	v.initDistances()
	v.initDiagram()
	v.initSeeds()
	v.initTessellation()
}

// TODO
func (v *Voronoi) initDistances() {

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

// initTessellation prepares the
func (v *Voronoi) initTessellation() {

	v.radius = 0
	v.activeSeeds = v.seeds

	// fmt.Println("####################################")
	// fmt.Println("#### Voronoi computation inited ####")
	// fmt.Println("####################################")
}

func (v *Voronoi) Tessellate(hideIterations bool) error {

	for len(v.activeSeeds) > 0 {

		stillActiveSeeds := []Point{}
		proximityDistances := v.getProximityDistances()

		for _, seed := range v.activeSeeds {
			// fmt.Println("Iteration starting. Active seeds: ", len(v.activeSeeds))

			stillActive := false
			for _, distance := range proximityDistances {
				stillActive = v.assignPointsToSeed(seed, v.distances[distance.X][distance.Y], distance.X, distance.Y) || stillActive
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

// assignPointsToSeed assigns all the eligible points at a given distance to the seed
func (v *Voronoi) assignPointsToSeed(seed Point, distance int, dx int, dy int) bool {
	stillActive := false
	stillActive = v.assignPointToSeed(seed, distance, dx, dy) || stillActive
	stillActive = v.assignPointToSeed(seed, distance, dx, -dy) || stillActive
	stillActive = v.assignPointToSeed(seed, distance, dy, dx) || stillActive
	stillActive = v.assignPointToSeed(seed, distance, dy, -dx) || stillActive
	stillActive = v.assignPointToSeed(seed, distance, -dx, dy) || stillActive
	stillActive = v.assignPointToSeed(seed, distance, -dx, -dy) || stillActive
	stillActive = v.assignPointToSeed(seed, distance, -dy, dx) || stillActive
	stillActive = v.assignPointToSeed(seed, distance, -dy, -dx) || stillActive

	return stillActive
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

func (v *Voronoi) getProximityDistances() []Point {
	combinations := []Point{}

	v.radius++
	dx := 0
	dy := v.radius

	for dy >= dx {
		combinations = append(combinations, Point{
			X: dx,
			Y: dy,
		})
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
