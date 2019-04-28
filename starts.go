package main

import (
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	starSize = 5
)

type star struct {
	renderer *sdl.Renderer
	tex      *sdl.Texture
	x, y     float64
	speed    float64
	texXPos  int32
	active   bool
}

func newStar(renderer *sdl.Renderer) *star {

	s := star{
		renderer: renderer,
		tex:      newTexture(renderer, "sprites/stars.png")}

	return &s
}

func (s *star) draw() {

	if !s.active {
		return
	}

	drawSize := starSize * scale
	drawX := s.x - float64(drawSize)/2.0
	drawY := s.y - float64(drawSize)/2.0
	drawRect := sdl.Rect{X: int32(drawX), Y: int32(drawY), W: int32(drawSize), H: int32(drawSize)}

	xTex := s.texXPos * starSize

	s.renderer.CopyEx(s.tex,
		&sdl.Rect{X: xTex, Y: 0, W: starSize, H: starSize},
		&drawRect,
		0,
		&sdl.Point{X: drawRect.W / 2.0, Y: drawRect.H / 2.0},
		sdl.FLIP_NONE)

	//s.renderer.SetDrawColor(255, 255, 255, 255)
	//s.renderer.DrawPoint(int32(s.x), int32(s.y))
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

	s.texXPos = int32(rand.Intn(7))
	s.x = float64(rand.Intn(int(screenWidth)))
	s.y = -10
	s.speed = 3 + rand.Float64()*3
	s.active = true
	//fmt.Printf("star start %v\n", s)
}

var stars []*star

func initStars(renderer *sdl.Renderer) {
	for i := 0; i < 300; i++ {
		s := newStar(renderer)
		stars = append(stars, s)
	}
	deltaStar = rand.Float64() * 8
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
		deltaStar = rand.Float64() * 8
	}
}
