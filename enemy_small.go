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
	tex                *sdl.Texture
	initX, x, y, angle float64
	active             bool
	texXPos            int32
	lastTimeBooster    time.Time
	index              int8
	turnLeft           bool
	fase               int8
}

func newEnemySmall(renderer *sdl.Renderer, x, y float64) (e enemySmall) {

	e.tex = newTexture(renderer, "sprites/enemy-small.png")
	e.x = x
	e.initX = x
	e.y = y
	e.angle = -90
	e.active = false

	e.texXPos = 0

	return e
}

func (e *enemySmall) draw(renderer *sdl.Renderer) {

	if !e.active {
		return
	}

	drawEnemySize := enemySmallSize * scale
	drawX := e.x - float64(drawEnemySize)/2.0
	drawY := e.y - float64(drawEnemySize)/2.0
	drawRect := sdl.Rect{X: int32(drawX), Y: int32(drawY), W: int32(drawEnemySize), H: int32(drawEnemySize)}

	xTex := e.texXPos * enemySmallSize

	renderer.CopyEx(e.tex,
		&sdl.Rect{X: xTex, Y: 0, W: enemySmallSize, H: enemySmallSize},
		&drawRect,
		(e.angle*-1)-90,
		&sdl.Point{X: drawRect.W / 2.0, Y: drawRect.H / 2.0},
		sdl.FLIP_NONE)

	if deb.active {
		// debug drawRect
		renderer.SetDrawColor(255, 0, 0, 255)
		renderer.DrawRect(&drawRect)

		// debug rect of collision
		collisionRect := sdl.Rect{X: int32(e.x - enemySmallCollisionRadius), Y: int32(e.y - enemySmallCollisionRadius),
			W: enemySmallCollisionRadius * 2, H: enemySmallCollisionRadius * 2}
		renderer.DrawRect(&collisionRect)

		// debug e.initX
		renderer.SetDrawColor(0, 0, 255, 255)
		renderer.DrawLine(int32(e.initX), 0, int32(e.initX), screenHeight)
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

func (e *enemySmall) start(index, fase int8, x, y float64) {

	e.index = index
	e.fase = fase
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

var enemiesSmall []*enemySmall

func initEnemiesSmall(renderer *sdl.Renderer) {
	for i := 0; i < 50; i++ {
		en := newEnemySmall(renderer, screenWidth/2+enemySmallSize, -1*enemySmallSize)
		enemiesSmall = append(enemiesSmall, &en)
	}
}

func deactivateAllEnemiesSmall() {
	for _, en := range enemiesSmall {
		en.active = false
	}
}

func enemySmallFromPool() *enemySmall {
	for _, en := range enemiesSmall {
		if !en.active {
			return en
		}
	}

	return nil
}
