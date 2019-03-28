package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var font *ttf.Font

func initFonts() {

	var err error

	fmt.Println("Init fonts")

	if err = ttf.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize TTF: %s\n", err)
		return
	}

	if font, err = ttf.OpenFont("fonts/m6x11.ttf", 18); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open font: %s\n", err)
		return
	}
}

func newText(renderer *sdl.Renderer, text string) (tex *sdl.Texture) {

	color := sdl.Color{R: 200, G: 200, B: 200, A: 255}

	var solid *sdl.Surface
	var err error
	if solid, err = font.RenderUTF8Solid(text, color); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to render text: %s\n", err)
		return
	}
	defer solid.Free()

	tex, err = renderer.CreateTextureFromSurface(solid)
	if err != nil {
		panic(fmt.Errorf("creating texture from %v: %v", text, err))
	}

	return tex
}
