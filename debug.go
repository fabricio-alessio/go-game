package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type debug struct {
	texs     [8]*sdl.Texture
	x, y     float64
	renderer *sdl.Renderer
	values   [8]string
	active   bool
}

func newDebug(renderer *sdl.Renderer) (d debug) {

	d.renderer = renderer
	d.y = 100

	return d
}

func (d *debug) draw(renderer *sdl.Renderer) {

	if !d.active {
		return
	}

	for pos, tex := range d.texs {
		d.drawTexture(pos, tex, renderer)
	}
}

func (d *debug) drawTexture(pos int, tex *sdl.Texture, renderer *sdl.Renderer) {

	if tex == nil {
		return
	}
	_, _, texW, texH, err := tex.Query()
	if err != nil {
		panic("Impossible to get texture dimentions from debug")
	}

	delta := int32(pos)*texH + 10
	y := d.y + float64(delta)
	renderer.Copy(tex,
		&sdl.Rect{X: 0, Y: 0, W: texW, H: texH},
		&sdl.Rect{X: int32(d.x), Y: int32(y), W: texW, H: texH})
}

func (d *debug) set(pos int8, value string) {

	if pos < 0 || pos > 5 {
		return
	}

	d.values[pos] = value
	d.update(pos)
}

func (d *debug) update(pos int8) {

	if !d.active {
		return
	}
	output := d.values[pos]
	d.texs[pos].Destroy()
	d.texs[pos] = newTextTexture(d.renderer, output)
}
