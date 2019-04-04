package main

import (
	"fmt"
	"os"

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

var startTick uint32

func startFPSTick() {

	startTick = sdl.GetTicks()
}

func stickWithFPS(fps uint32) {

	deltaTick := sdl.GetTicks() - startTick
	if 1000/fps > deltaTick {
		sdl.Delay(1000/fps - deltaTick)
	}
}

func initSdl() {

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDL: %s\n", err)
		panic(err)
	}
}

func initScreen(title string, width, height int32, fullScreen bool) (*sdl.Window, *sdl.Renderer) {

	var window *sdl.Window
	var renderer *sdl.Renderer
	var err error

	if fullScreen {
		window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
			width, height, sdl.WINDOW_FULLSCREEN_DESKTOP)
	} else {
		window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
			width, height, sdl.WINDOW_OPENGL)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		panic(err)
	}

	sdl.ShowCursor(sdl.DISABLE)

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		panic(err)
	}

	return window, renderer
}

func showScreenInfo(renderer *sdl.Renderer) {

	w, h, err := renderer.GetOutputSize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't get output size: %s\n", err)
		panic(err)
	} else {
		fmt.Println(w, h)
	}
}

var joystic0 *sdl.Joystick

func initJoystic(index int) {

	fmt.Fprintf(os.Stdout, "NumJoysticks: %v\n", sdl.NumJoysticks())
	if sdl.NumJoysticks() < 1 {
		return
	}
	joy := sdl.JoystickOpen(index)
	if joy != nil {
		fmt.Fprintf(os.Stdout, "Name: %s\n", sdl.JoystickNameForIndex(0))
		fmt.Fprintf(os.Stdout, "Number of Axes: %d\n", joy.NumAxes())
		fmt.Fprintf(os.Stdout, "Number of Buttons: %d\n", joy.NumButtons())
		fmt.Fprintf(os.Stdout, "Number of Balls: %d\n", joy.NumBalls())
	} else {
		fmt.Fprintf(os.Stderr, "Couldn't open Joystick 0\n")
	}

	if index == 0 {
		joystic0 = joy
	}

	// Close if opened
	//	if joy.Attached() {
	//		joy.Close()
	//	}
}

func showJoysticState(index int) {

	var joy *sdl.Joystick
	if index == 0 {
		joy = joystic0
	}

	for b := 0; b < joy.NumButtons(); b++ {
		if joy.Button(b) > 0 {
			fmt.Fprintf(os.Stdout, "Button %d: %d\n", b, joy.Button(b))
		}
	}
	//for a := 0; a < joy.NumAxes(); a++ {
	//	if joy.Axis(a) > 0 {
	//		fmt.Fprintf(os.Stdout, "Axis %d: %d\n", a, joy.Axis(a))
	//	}
	//}

	fmt.Fprintf(os.Stdout, "Axis %d: %d\n", 0, joy.Axis(0))
}
