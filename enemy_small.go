package main

import (
	"fmt"
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	enemySmallSpeed           = 2
	enemySmallSize            = 16
	enemySmallCollisionRadius = 16
	enemySmallRadius          = 30
)

type enemySmall struct {
	renderer           *sdl.Renderer
	tex                *sdl.Texture
	initX, x, y, angle float64
	active             bool
	texXPos            int32
	lastTimeBooster    time.Time
	index              int8
	turnLeft           bool
	fase               int8
}

func newEnemySmall(renderer *sdl.Renderer, x, y float64) *enemySmall {

	enemy := enemySmall{
		renderer: renderer,
		tex:      newTexture(renderer, "sprites/enemy-small.png"),
		x:        x,
		initX:    x,
		y:        y,
		angle:    -90,
		active:   false,
		texXPos:  0}

	return &enemy
}

func (e *enemySmall) draw() {

	if !e.active {
		return
	}

	drawEnemySize := enemySmallSize * scale
	drawX := e.x - float64(drawEnemySize)/2.0
	drawY := e.y - float64(drawEnemySize)/2.0
	drawRect := sdl.Rect{X: int32(drawX), Y: int32(drawY), W: int32(drawEnemySize), H: int32(drawEnemySize)}

	xTex := e.texXPos * enemySmallSize

	e.renderer.CopyEx(e.tex,
		&sdl.Rect{X: xTex, Y: 0, W: enemySmallSize, H: enemySmallSize},
		&drawRect,
		(e.angle*-1)-90,
		&sdl.Point{X: drawRect.W / 2.0, Y: drawRect.H / 2.0},
		sdl.FLIP_NONE)

	if deb.active {

		//fmt.Printf("%v\n", drawRect)
		// debug drawRect
		e.renderer.SetDrawColor(255, 0, 0, 255)
		e.renderer.DrawRect(&drawRect)

		// debug rect of collision
		collisionRect := sdl.Rect{X: int32(e.x - enemySmallCollisionRadius), Y: int32(e.y - enemySmallCollisionRadius),
			W: enemySmallCollisionRadius * 2, H: enemySmallCollisionRadius * 2}
		e.renderer.DrawRect(&collisionRect)

		// debug e.initX
		e.renderer.SetDrawColor(0, 0, 255, 255)
		e.renderer.DrawLine(int32(e.initX), 0, int32(e.initX), screenHeight)
	}
}

func (e *enemySmall) update() {

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

	eSpeed := enemySmallSpeed * 2 * delta

	if e.turnLeft {
		e.angle += eSpeed / 2
		if e.angle >= 45 {
			e.turnLeft = false
		}
	} else {
		e.angle -= eSpeed / 2
		if e.angle <= -225 {
			e.turnLeft = true
		}
	}

	radiansAngle := e.angle * math.Pi / 180
	xSpeed := 5 * delta * math.Cos(radiansAngle)
	ySpeed := 5 * delta * math.Sin(radiansAngle*-1)

	e.x += xSpeed
	e.y += ySpeed
	if e.y > screenHeight {
		e.active = false
	}

	if e.index == 0 {
		deb.set(2, fmt.Sprintf("angle %f", e.angle))
		deb.set(3, fmt.Sprintf("e.x %f", e.x))
		deb.set(4, fmt.Sprintf("e.y %f", e.y))
	}
}

func (e *enemySmall) start(x, y, angle, speed float64, entityType int8) {

	e.angle = 0
	e.x = x
	e.initX = x
	e.y = y
	e.active = true
}

func (e *enemySmall) beHit() {

	e.beDestroyed()
}

func (e *enemySmall) beDestroyed() {

	score.incrementPointsP1(1)
	chunkExplosion.Play(2, 0)
	e.active = false
	ex := explosionFromPool()
	ex.start(e.x, e.y, enemySmallSpeed)
}

func (e *enemySmall) executeCollisionWith(other entity) {

	if other.getType() == entityTypePlayerBullet {
		e.beDestroyed()
	} else if other.getType() == entityTypePlayer {
		e.beDestroyed()
	}
}

func (e *enemySmall) getCollisionCircle() circle {

	return circle{x: e.x, y: e.y, radius: enemySmallCollisionRadius}
}

func (e *enemySmall) isActive() bool {

	return e.active
}

func (e *enemySmall) getType() int8 {

	return entityTypeEnemySmall
}

func (e *enemySmall) deactivate() {

	e.active = false
}

var enemiesSmall []entity

func initEnemiesSmall(renderer *sdl.Renderer) {
	for i := 0; i < 80; i++ {
		en := newEnemySmall(renderer, screenWidth/2+enemySmallSize, -1*enemySmallSize)
		enemiesSmall = append(enemiesSmall, en)
	}
}

func deactivateAllEnemiesSmall() {
	for _, en := range enemiesSmall {
		en.deactivate()
	}
}

func enemySmallFromPool() entity {
	for _, en := range enemiesSmall {
		if !en.isActive() {
			return en
		}
	}

	return nil
}
