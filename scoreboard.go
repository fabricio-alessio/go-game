package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type scoreboard struct {
	tex      *sdl.Texture
	x, y     float64
	points   int32
	renderer *sdl.Renderer
}

func newScoreboard(renderer *sdl.Renderer) (s scoreboard) {

	s.points = 0
	s.renderer = renderer
	s.update()

	return s
}

func (s *scoreboard) draw(renderer *sdl.Renderer) {

	_, _, texW, texH, err := s.tex.Query()
	if err != nil {
		panic("Impossible to get texture dimentions from scoreboard")
	}

	renderer.Copy(s.tex,
		&sdl.Rect{X: 0, Y: 0, W: texW, H: texH},
		&sdl.Rect{X: int32(s.x), Y: int32(s.y), W: texW * scale, H: texH * scale})
}

func (s *scoreboard) reset() {

	s.points = 0
	s.update()
}

func (s *scoreboard) increment(qtd int32) {

	s.points += qtd
	s.update()
}

func (s *scoreboard) update() {

	output := fmt.Sprintf("%d", s.points)
	s.tex.Destroy()
	s.tex = newText(s.renderer, output)
}
