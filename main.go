package main

import (
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	//screenWidth  = 640
	//screenHeight = 480
	screenWidth  = 1536
	screenHeight = 864
	//screenWidth          = 800
	//screenHeight         = 600
	targetTicksPerSecond = 60
	millisByFrame        = (uint32)(1000 / targetTicksPerSecond)
	scale                = 2
)

var winTitle string = "Game"
var delta float64
var joy *sdl.Joystick
var score scoreboard
var deb debug
var plr player
var efleet enemyFleet
var gameStarted bool

func joysticsDetails() {
	fmt.Fprintf(os.Stdout, "NumJoysticks: %v\n", sdl.NumJoysticks())
	if sdl.NumJoysticks() < 1 {
		return
	}
	joy = sdl.JoystickOpen(0)
	if joy != nil {
		fmt.Fprintf(os.Stdout, "Name: %s\n", sdl.JoystickNameForIndex(0))
		fmt.Fprintf(os.Stdout, "Number of Axes: %d\n", joy.NumAxes())
		fmt.Fprintf(os.Stdout, "Number of Buttons: %d\n", joy.NumButtons())
		fmt.Fprintf(os.Stdout, "Number of Balls: %d\n", joy.NumBalls())
	} else {
		fmt.Fprintf(os.Stderr, "Couldn't open Joystick 0\n")
	}

	// Close if opened
	//	if joy.Attached() {
	//		joy.Close()
	//	}
}

func joysticState() {
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

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDL: %s\n", err)
		panic(err)
	}

	joysticsDetails()
	initSounds()
	initFonts()

	var window *sdl.Window
	var renderer *sdl.Renderer
	var err error

	window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWidth, screenHeight, sdl.WINDOW_FULLSCREEN_DESKTOP)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		panic(err)
	}
	defer window.Destroy()

	sdl.ShowCursor(sdl.DISABLE)

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		panic(err)
	}
	defer renderer.Destroy()

	gameMenu := newMenu(renderer)

	score = newScoreboard(renderer)
	deb = newDebug(renderer)

	plr = newPlayer(renderer)

	initEnemies(renderer)
	initExplosions(renderer)
	initBulletPool(renderer)
	efleet = newEnemyFleet()
	efleet.startLevel(1)

	plr.start()

	w, h, err := renderer.GetOutputSize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't get output size: %s\n", err)
		panic(err)
	} else {
		fmt.Println(w, h)
	}

	gameStarted = false

	paused := false
	lastTimePauseToggle := time.Now()
	lastTimeDebugToggle := time.Now()
	for {
		frameStartTime := time.Now()
		sdl.PumpEvents()
		keys := sdl.GetKeyboardState()
		if keys[sdl.SCANCODE_ESCAPE] == 1 {
			return
		}

		renderer.SetDrawColor(10, 10, 10, 255)
		renderer.Clear()

		if gameStarted {

			if keys[sdl.SCANCODE_F1] == 1 {
				if time.Since(lastTimePauseToggle) >= playerShotCooldown {
					paused = !paused
					lastTimePauseToggle = time.Now()
				}
			}
			if keys[sdl.SCANCODE_F2] == 1 {
				if time.Since(lastTimeDebugToggle) >= playerShotCooldown {
					deb.active = !deb.active
					lastTimeDebugToggle = time.Now()
				}
			}

			if paused {
				sdl.Delay(100)
				continue
			}

			plr.draw(renderer)
			plr.update()
			efleet.update()
			for _, en := range enemies {
				en.draw(renderer)
				en.update()
			}
			for _, ex := range explosions {
				ex.draw(renderer)
				ex.update()
			}

			for _, bul := range bulletPool {
				bul.draw(renderer)
				bul.update()
			}
			score.draw(renderer)
			deb.draw(renderer)
			checkCollisions(enemies)
		} else {

			gameMenu.draw(renderer)
			gameMenu.update()
		}

		nanosSince := time.Since(frameStartTime).Nanoseconds()
		millisSince := (uint32)(nanosSince / 1000000)
		delay := millisByFrame - (millisSince + 6)
		if delay > 0 {
			sdl.Delay(delay)
		}
		renderer.Present()

		delta = time.Since(frameStartTime).Seconds() * targetTicksPerSecond
	}
}
