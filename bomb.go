package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	bombSpeed           = 2
	bombSize            = 16
	bombCollisionRadius = 16
)

type bomb struct {
	renderer        *sdl.Renderer
	tex             *sdl.Texture
	x, y            float64
	active          bool
	texXPos         int32
	lastTimeBooster time.Time
}

func newBomb(renderer *sdl.Renderer, x, y float64) *bomb {

	b := bomb{
		renderer: renderer,
		tex:      newTexture(renderer, "sprites/bomb.png"),
		x:        x,
		y:        y,
		active:   false,
		texXPos:  0}

	return &b
}

func (e *bomb) draw() {

	if !e.active {
		return
	}

	drawSize := bombSize * scale
	drawX := e.x - float64(drawSize)/2.0
	drawY := e.y - float64(drawSize)/2.0
	drawRect := sdl.Rect{X: int32(drawX), Y: int32(drawY), W: int32(drawSize), H: int32(drawSize)}

	xTex := e.texXPos * bombSize

	e.renderer.CopyEx(e.tex,
		&sdl.Rect{X: xTex, Y: 0, W: bombSize, H: bombSize},
		&drawRect,
		0,
		&sdl.Point{X: drawRect.W / 2.0, Y: drawRect.H / 2.0},
		sdl.FLIP_NONE)

	if deb.active {

		// debug drawRect
		e.renderer.SetDrawColor(255, 0, 0, 255)
		e.renderer.DrawRect(&drawRect)

		// debug rect of collision
		collisionRect := sdl.Rect{X: int32(e.x - bombCollisionRadius), Y: int32(e.y - bombCollisionRadius),
			W: bombCollisionRadius * 2, H: bombCollisionRadius * 2}
		e.renderer.DrawRect(&collisionRect)
	}
}

func (e *bomb) update() {

	if !e.active {
		return
	}

	if time.Since(e.lastTimeBooster) >= playerShotCooldown/2 {
		if e.texXPos == 0 {
			e.texXPos = 1
		} else if e.texXPos == 1 {
			e.texXPos = 0
		}
		e.lastTimeBooster = time.Now()
	}

	e.y += bombSpeed * delta
	if e.y > screenHeight {
		fmt.Println("Send bomb to poll")
		e.active = false
	}
}

func (e *bomb) start(x, y, angle, speed float64, entityType int8) {

	fmt.Println("Bomb start")
	e.x = x
	e.y = -30
	e.active = true
}

func (e *bomb) beDestroyed() {

	e.deactivate()
	mixer.playSound("explosion")
	ex := explosionFromPool()
	ex.start(e.x, e.y, enemySmallSpeed)
}

func (e *bomb) executeCollisionWith(other entity) {

	if other.getType() == entityTypePlayer {
		e.beDestroyed()
	}
}

func (e *bomb) getCollisionCircle() circle {

	return circle{x: e.x, y: e.y, radius: bombCollisionRadius}
}

func (e *bomb) isActive() bool {

	return e.active
}

func (e *bomb) getType() int8 {

	return entityTypeBomb
}

func (e *bomb) deactivate() {

	e.active = false
}

var bombs []entity

func initBombs(renderer *sdl.Renderer) {
	for i := 0; i < 30; i++ {
		en := newBomb(renderer, screenWidth/2+bombSize, -1*bombSize)
		bombs = append(bombs, en)
	}
}

func bombFromPool() entity {
	for _, en := range bombs {
		if !en.isActive() {
			return en
		}
	}

	return nil
}
