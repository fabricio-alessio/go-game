package main

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	//screenWidth  = 640
	//screenHeight = 480
	//screenWidth  = 1536
	//screenHeight = 864
	screenWidth  = 800
	screenHeight = 600
	scale        = 2
	fps          = 60
)

var delta float64
var score scoreboard
var deb debug
var plr entity
var efleet enemyFleet
var gameStarted bool
var mixer *soundManager

func main() {

	initSdl()
	initJoystic(0)
	initSounds()
	mixer = newSoundManager()
	mixer.loadSound("hit", "sounds/hit.wav")
	mixer.loadSound("laser", "sounds/laser.wav")
	//mixer.loadSound("enemyLaser", "sounds/laser.wav")
	mixer.loadSound("explosion", "sounds/explosion.wav")
	initFonts()
	var window, renderer = initScreen("Game", screenWidth, screenHeight, false)
	defer window.Destroy()
	defer renderer.Destroy()

	gameMenu := newMenu(renderer)
	score = newScoreboard(renderer)
	deb = newDebug(renderer)
	plr = newPlayer(renderer)

	initEnemiesBig(renderer)
	initEnemiesSmall(renderer)
	initExplosions(renderer)
	initBulletPool(renderer)
	efleet = newEnemyFleet()
	efleet.startLevel(1)

	plr.start(0, 0, 0, 0, 0)

	gameStarted = false

	paused := false
	lastTimePauseToggle := time.Now()
	lastTimeDebugToggle := time.Now()
	for {
		startFPSTick()

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

			plr.draw()
			plr.update()
			efleet.update()
			for _, en := range enemiesBig {
				en.draw()
				en.update()
			}
			for _, en := range enemiesSmall {
				en.draw()
				en.update()
			}
			for _, ex := range explosions {
				ex.draw(renderer)
				ex.update()
			}

			for _, bul := range bulletPool {
				bul.draw()
				bul.update()
			}
			score.draw(renderer)
			deb.draw(renderer)
			checkCollisions(enemiesBig, enemiesSmall)
		} else {

			gameMenu.draw(renderer)
			gameMenu.update()
		}

		renderer.Present()
		stickWithFPS(fps)

		delta = time.Since(frameStartTime).Seconds() * fps
	}
}
