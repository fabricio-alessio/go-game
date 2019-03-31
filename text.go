package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type text struct {
	tex        *sdl.Texture
	x, y       float64
	value      string
	renderer   *sdl.Renderer
	texW, texH int32
}

func newText(renderer *sdl.Renderer, x, y float64, value string) (t text) {

	t.renderer = renderer
	t.x = x
	t.y = y
	t.setValue(value)

	return t
}

func (t *text) draw(renderer *sdl.Renderer) {

	renderer.Copy(t.tex,
		&sdl.Rect{X: 0, Y: 0, W: t.texW, H: t.texH},
		&sdl.Rect{X: int32(t.x), Y: int32(t.y), W: t.texW * scale, H: t.texH * scale})
}

func (t *text) setValue(newValue string) {

	t.value = newValue
	t.tex.Destroy()
	t.tex = newTextTexture(t.renderer, newValue)

	var err error
	_, _, t.texW, t.texH, err = t.tex.Query()
	if err != nil {
		panic("Impossible to get texture dimentions from text")
	}
}

func (t *text) update() {

	// do nothing
}
