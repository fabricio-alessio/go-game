package main

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	menuChoice1P = 1
	menuChoice2P = 2
)

type menu struct {
	tex             *sdl.Texture
	navyAngle       float64
	navyZoom        int8
	navyX, navyY    float64
	navySpeedX      float64
	navyDirection   float64
	navyInclination int32
	texYPos         int32
	choice          int8
	lastTimeBooster time.Time
}

func newMenu(renderer *sdl.Renderer) (m menu) {

	m.tex = newTexture(renderer, "sprites/ship.png")
	m.choice = menuChoice1P
	m.navyX = screenWidth / 2
	m.navyY = screenHeight / 3
	m.navyAngle = 0
	m.navyZoom = 5
	m.navySpeedX = 5
	m.navyDirection = 1
	m.navyInclination = 2
	m.texYPos = 0
	return m
}

func (m *menu) draw(renderer *sdl.Renderer) {

	drawPlayerWidth := playerWidth * m.navyZoom
	drawPlayerHeight := playerHeight * m.navyZoom
	drawX := m.navyX - float64(drawPlayerWidth)/2.0
	drawY := m.navyY - float64(drawPlayerHeight)/2.0
	drawRect := sdl.Rect{X: int32(drawX), Y: int32(drawY), W: int32(drawPlayerWidth), H: int32(drawPlayerHeight)}

	xTex := int32(m.navyInclination * playerWidth)
	yTex := int32(m.texYPos * playerHeight)

	renderer.CopyEx(m.tex,
		&sdl.Rect{X: xTex, Y: yTex, W: playerWidth, H: playerHeight},
		&drawRect,
		0,
		&sdl.Point{X: drawRect.W / 2.0, Y: drawRect.H / 2.0},
		sdl.FLIP_NONE)
}

func (m *menu) update() {

	m.handleKeyboard()
	m.handleJoystic()

	m.navyX += m.navySpeedX * m.navyDirection * delta

	border := float64(screenWidth / 4)
	if m.navyX < border {
		m.navyDirection = 1
	} else if m.navyX > screenWidth-border {
		m.navyDirection = -1
	}

	step := (border * 2) / 5
	if m.navyX > border && m.navyX < border+step {
		m.navyInclination = 0
		m.navySpeedX = 3.5
	} else if m.navyX > border+step && m.navyX < border+(2*step) {
		m.navyInclination = 1
		m.navySpeedX = 4
	} else if m.navyX > border+(2*step) && m.navyX < border+(3*step) {
		m.navyInclination = 2
		m.navySpeedX = 4.5
	} else if m.navyX > border+(3*step) && m.navyX < border+(4*step) {
		m.navyInclination = 3
		m.navySpeedX = 4
	} else if m.navyX > border+(4*step) && m.navyX < border+(5*step) {
		m.navyInclination = 4
		m.navySpeedX = 3.5
	}

	if time.Since(m.lastTimeBooster) >= playerShotCooldown/2 {
		if m.texYPos == 0 {
			m.texYPos = 1
		} else if m.texYPos == 1 {
			m.texYPos = 0
		}
		m.lastTimeBooster = time.Now()
	}
}

func (m *menu) handleKeyboard() {

	keys := sdl.GetKeyboardState()

	if keys[sdl.SCANCODE_UP] == 1 {
		m.moveUp()
	} else if keys[sdl.SCANCODE_DOWN] == 1 {
		m.moveDown()
	}

	if keys[sdl.SCANCODE_SPACE] == 1 {
		m.choose()
	}
}

func (m *menu) handleJoystic() {

	if joy.Axis(1) < -8000 {
		m.moveUp()
	} else if joy.Axis(1) > 8000 {
		m.moveDown()
	}

	if joy.Button(1) > 0 {
		m.choose()
	}
}

func (m *menu) moveUp() {

	if m.choice == menuChoice1P {
		m.choice = menuChoice2P
	} else {
		m.choice = menuChoice1P
	}
}

func (m *menu) moveDown() {

	if m.choice == menuChoice1P {
		m.choice = menuChoice2P
	} else {
		m.choice = menuChoice1P
	}
}

func (m *menu) choose() {

	gameStarted = true
	score.resetPointsP1()
	plr.replay()
	plr.start()
}
