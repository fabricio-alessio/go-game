package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	powerUpSpeed           = 2
	powerUpSize            = 16
	powerUpCollisionRadius = 16
)

type powerUp struct {
	renderer        *sdl.Renderer
	tex             *sdl.Texture
	x, y            float64
	active          bool
	texXPos         int32
	lastTimeBooster time.Time
	entityType      int8
}

func newPowerUp(renderer *sdl.Renderer, x, y float64) *powerUp {

	power := powerUp{
		renderer: renderer,
		tex:      newTexture(renderer, "sprites/power-up.png"),
		x:        x,
		y:        y,
		active:   false,
		texXPos:  0}

	return &power
}

func (e *powerUp) draw() {

	if !e.active {
		return
	}

	drawPowerUpSize := powerUpSize * scale
	drawX := e.x - float64(drawPowerUpSize)/2.0
	drawY := e.y - float64(drawPowerUpSize)/2.0
	drawRect := sdl.Rect{X: int32(drawX), Y: int32(drawY), W: int32(drawPowerUpSize), H: int32(drawPowerUpSize)}

	xTex := e.texXPos * powerUpSize
	var yTex int32
	if e.entityType == entityTypePowerUpBullet {
		yTex = 0
	} else {
		yTex = powerUpSize
	}

	e.renderer.CopyEx(e.tex,
		&sdl.Rect{X: xTex, Y: yTex, W: powerUpSize, H: powerUpSize},
		&drawRect,
		0,
		&sdl.Point{X: drawRect.W / 2.0, Y: drawRect.H / 2.0},
		sdl.FLIP_NONE)

	if deb.active {

		// debug drawRect
		e.renderer.SetDrawColor(255, 0, 0, 255)
		e.renderer.DrawRect(&drawRect)

		// debug rect of collision
		collisionRect := sdl.Rect{X: int32(e.x - powerUpCollisionRadius), Y: int32(e.y - powerUpCollisionRadius),
			W: powerUpCollisionRadius * 2, H: powerUpCollisionRadius * 2}
		e.renderer.DrawRect(&collisionRect)
	}
}

func (e *powerUp) update() {

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

	e.y += enemySmallSpeed * delta
	if e.y > screenHeight {
		fmt.Println("Send power up to poll")
		e.active = false
	}
}

func (e *powerUp) start(x, y, angle, speed float64, entityType int8) {

	fmt.Println("Power up start")
	e.x = x
	e.y = -30
	e.active = true
	e.entityType = entityType
}

func (e *powerUp) beCollected() {

	e.deactivate()
	mixer.playSound("powerUp")
}

func (e *powerUp) executeCollisionWith(other entity) {

	if other.getType() == entityTypePlayer {
		e.beCollected()
	}
}

func (e *powerUp) getCollisionCircle() circle {

	return circle{x: e.x, y: e.y, radius: powerUpCollisionRadius}
}

func (e *powerUp) isActive() bool {

	return e.active
}

func (e *powerUp) getType() int8 {

	return e.entityType
}

func (e *powerUp) deactivate() {

	e.active = false
}

var powerUps []entity

func initPowerUps(renderer *sdl.Renderer) {
	for i := 0; i < 30; i++ {
		en := newPowerUp(renderer, screenWidth/2+powerUpSize, -1*powerUpSize)
		powerUps = append(powerUps, en)
	}
}

func powerUpFromPool() entity {
	for _, en := range powerUps {
		if !en.isActive() {
			return en
		}
	}

	return nil
}
