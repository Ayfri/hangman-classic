package main

import "math"

type Angle struct {
	Degrees float64
    Radians float64
}

type Vector2 struct {
	X, Y, length float64
	angle Angle
}

// Collides Create a function to test if two vectors will collide
func (v1 Vector2) Collides(v2 Vector2) bool {
    // Check if the distance between the two vectors is less than the sum of their radii
    if v1.length + v2.length < v1.distance(v2) {
        return true
    }
    return false
}

func (v1 Vector2) distance(v2 Vector2) float64 {
    // Calculate the distance between the two vectors
    return math.Sqrt(math.Pow(v1.X - v2.X, 2) + math.Pow(v1.Y - v2.Y, 2))
}

type Box struct {
	X, Y, Width, Height float64
}

// test if two boxes collides
func (b1 Box) Collides(b2 Box) bool {
    // Check if the distance between the two vectors is less than the sum of their radii
    if b1.X + b1.Width > b2.X && b1.X < b2.X + b2.Width && b1.Y + b1.Height > b2.Y && b1.Y < b2.Y + b2.Height {
        return true
    }
    return false
}
