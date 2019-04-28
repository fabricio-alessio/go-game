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
	enemySmallValCount        = 4
	enemySmallValH            = 8
)

var enemyVariations = [enemySmallValCount][enemySmallValH]float64{
	{0, 1, -1, -1, 1, 1, 0, 0},
	{0, -1, 1, 1, -1, -1, 0, 0},
	{0, 1, 0, -1, 0, 1, 0, -1},
	{0, -1, 0, 1, 0, -1, 0, 1}}

type enemySmall struct {
	renderer        *sdl.Renderer
	tex             *sdl.Texture
	x, y, angle     float64
	active          bool
	texXPos         int32
	lastTimeBooster time.Time
	group           *enemyGroup
}

func newEnemySmall(renderer *sdl.Renderer, x, y float64) *enemySmall {

	enemy := enemySmall{
		renderer: renderer,
		tex:      newTexture(renderer, "sprites/enemy-small.png"),
		x:        x,
		y:        y,
		angle:    270,
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

		// debug drawRect
		e.renderer.SetDrawColor(255, 0, 0, 255)
		e.renderer.DrawRect(&drawRect)

		// debug rect of collision
		collisionRect := sdl.Rect{X: int32(e.x - enemySmallCollisionRadius), Y: int32(e.y - enemySmallCollisionRadius),
			W: enemySmallCollisionRadius * 2, H: enemySmallCollisionRadius * 2}
		e.renderer.DrawRect(&collisionRect)
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

	step := screenHeight / enemySmallValH
	e.y += enemySmallSpeed * delta
	index := int(e.y / float64(step))

	if index >= 0 && index < enemySmallValH {
		e.angle += delta * enemyVariations[e.group.variationIndex][index]
	}

	radiansAngle := e.angle * math.Pi / 180
	e.x += 2 * math.Cos(radiansAngle)

	if e.y > screenHeight+enemySmallSize*scale {
		fmt.Println("Send small enemy back to poll")
		e.active = false
	}
}

func (e *enemySmall) start(x, y, angle, speed float64, entityType int8) {

	fmt.Println("Enemy small start")
	e.angle = 270
	e.x = x
	e.y = y
	e.active = true
}

func (e *enemySmall) beHit() {

	e.beDestroyed()
}

func (e *enemySmall) beDestroyed() {

	score.incrementPointsP1(1)
	mixer.playSound("explosion")
	e.active = false
	e.group.aliveCount--
	if e.group.aliveCount == 0 {
		fmt.Println("All small group dead")
		pu := powerUpFromPool()
		pu.start(e.x, 0, 0, 0, entityTypePowerUpBullet)
	}
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
