package main

type Color struct {
	R byte
	G byte
	B byte
	A byte
}

type Point struct {
	X        int
	Y        int
	Distance *int
	Color    *Color
}
