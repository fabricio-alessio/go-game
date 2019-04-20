package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	enemyExtraCommingSpeed = 1
	enemyExtraCommingSize  = 22
)

type enemyExtraComming struct {
	renderer        *sdl.Renderer
	tex             *sdl.Texture
	x, y            float64
	active          bool
	lastTimeBooster time.Time
	texXPos         int32
}

func newEnemyExtraComming(renderer *sdl.Renderer, x, y float64) *enemyExtraComming {

	en := enemyExtraComming{
		renderer: renderer,
		tex:      newTexture(renderer, "sprites/enemy-extra-comming.png"),
		texXPos:  0,
		x:        x,
		y:        y,
		active:   false}

	return &en
}

func (e *enemyExtraComming) start(x, y, angle, speed float64, entityType int8) {

	e.x = float64(rand.Intn(int(screenWidth)))
	e.y = screenHeight + 30
	e.active = true
}

func (e *enemyExtraComming) executeCollisionWith(other entity) {

	// do nothing
}

func (e *enemyExtraComming) getCollisionCircle() circle {

	return circle{}
}

func (e *enemyExtraComming) isActive() bool {

	return e.active
}

func (e *enemyExtraComming) getType() int8 {

	return entityTypeEnemyExtraComming
}

func (e *enemyExtraComming) draw() {

	if !e.active {
		return
	}

	drawEnemySize := enemyExtraCommingSize * scale
	drawX := e.x - float64(drawEnemySize)/2.0
	drawY := e.y - float64(drawEnemySize)/2.0
	drawRect := sdl.Rect{X: int32(drawX), Y: int32(drawY), W: int32(drawEnemySize), H: int32(drawEnemySize)}

	xTex := e.texXPos * enemyExtraCommingSize

	e.renderer.CopyEx(e.tex,
		&sdl.Rect{X: xTex, Y: 0, W: enemyExtraCommingSize, H: enemyExtraCommingSize},
		&drawRect,
		0,
		&sdl.Point{X: drawRect.W / 2.0, Y: drawRect.H / 2.0},
		sdl.FLIP_NONE)

	if deb.active {
		// debug drawRect
		e.renderer.SetDrawColor(255, 0, 0, 255)
		e.renderer.DrawRect(&drawRect)
	}
}

func (e *enemyExtraComming) booster() {

	if time.Since(e.lastTimeBooster) >= playerShotCooldown/2 {
		if e.texXPos == 0 {
			e.texXPos = 1
		} else if e.texXPos == 1 {
			e.texXPos = 0
		}
		e.lastTimeBooster = time.Now()
	}
}

func (e *enemyExtraComming) update() {

	if !e.active {
		return
	}

	e.booster()

	e.y -= enemyExtraSpeed * delta
	if e.y < -30 {
		fmt.Println("Back to attack")
		e.active = false
		enemy := enemyExtraFromPool()
		enemy.start(e.x, 0, 0, 0, 0)
	}
}

func (e *enemyExtraComming) deactivate() {

	e.active = false
}

var enemiesExtraComming []entity

func initEnemiesExtraComming(renderer *sdl.Renderer) {
	for i := 0; i < 20; i++ {
		en := newEnemyExtraComming(renderer, screenWidth/2+enemyExtraCommingSize, -1*enemyExtraCommingSize)
		enemiesExtraComming = append(enemiesExtraComming, en)
	}
}

func deactivateAllEnemiesExtraComming() {
	for _, en := range enemiesExtraComming {
		en.deactivate()
	}
}

func enemyExtraCommingFromPool() entity {
	for _, en := range enemiesExtraComming {
		if !en.isActive() {
			return en
		}
	}

	return nil
}
