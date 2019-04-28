package main

import (
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

type star struct {
	renderer *sdl.Renderer
	x, y     float64
	speed    float64
	active   bool
}

func newStar(renderer *sdl.Renderer) *star {

	s := star{
		renderer: renderer}

	return &s
}

func (s *star) draw() {

	if !s.active {
		return
	}

	s.renderer.SetDrawColor(255, 255, 255, 255)
	s.renderer.DrawPoint(int32(s.x), int32(s.y))
	//fmt.Printf("draw star %v\n", s)
}

func (s *star) update() {

	if !s.active {
		return
	}

	s.y += s.speed * delta
	//fmt.Printf("update star %v\n", s)
	if s.y > screenHeight+5 {
		s.active = false
	}
}

func (s *star) start() {

	s.x = float64(rand.Intn(int(screenWidth)))
	s.y = -10
	s.speed = rand.Float64() * 5
	s.active = true
	//fmt.Printf("star start %v\n", s)
}

var stars []*star

func initStars(renderer *sdl.Renderer) {
	for i := 0; i < 300; i++ {
		s := newStar(renderer)
		stars = append(stars, s)
	}
	deltaStar = rand.Float64() * 4
}

func starFromPool() *star {
	for _, s := range stars {
		if !s.active {
			return s
		}
	}

	return nil
}

var space float64
var deltaStar float64

func releaseStars() {

	space += delta
	//fmt.Printf("space %v deltaStar %v\n", space, deltaStar)
	if space >= deltaStar {
		star := starFromPool()
		star.start()
		space = 0
		deltaStar = rand.Float64() * 4
	}
}
