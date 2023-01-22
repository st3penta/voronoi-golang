# voronoi-golang
Golang implementation of the voronoi diagram


## Overview
This is a graphical representation of a [voronoi diagram](https://en.wikipedia.org/wiki/Voronoi_diagram) written in Go.

Demo:  
<img alt="demo gif" src="demo.gif" width="100" height="100">


## Usage
Run the main script without any parameters.  
Some parameters may be customized [here](../main.go#L7-L21)


## Hotkeys

`Space`: restarts the simulation generating a new set of seeds  
`Enter`: starts/stops the simulation


## Something about the algorithm used
This solution implements an approximated (but enough accurate) algorithm to solve the voronoi diagram problem.  
Basically it works by extending the cells starting from the seed points: at each iteration, the size of the cell is incremented and each point laying in this incremental perimeter is assigned to the cell (unless it is already assigned to a closed seed).  
Each cell is extended until no more points can be assigned to it. At that point the cell is flagged as `inactive` and ignored in the following iterations.  
The algorithm goes on until there are no more active cells, meaning all the points in the canvas have been assigned to a cell.

**Why this kind of solution?**
Because it is in a sweet spot between the simple but highly inefficient brute force algorithm (that for each point in the canvas finds its nearest seed,) and the [Fortune's Algorithm](https://en.wikipedia.org/wiki/Fortune%27s_algorithm), efficient but complex