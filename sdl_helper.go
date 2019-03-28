package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func newTexture(renderer *sdl.Renderer, fileName string) (tex *sdl.Texture) {

	image, err := img.Load(fileName)
	if err != nil {
		panic(fmt.Errorf("loading file name %v: %v", fileName, err))
	}
	defer image.Free()
	tex, err = renderer.CreateTextureFromSurface(image)
	if err != nil {
		panic(fmt.Errorf("creating texture from %v: %v", fileName, err))
	}
	return tex
}
