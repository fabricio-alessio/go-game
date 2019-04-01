package main

import (
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	enemyBigSpeed           = 1
	enemyBigSize            = 32
	enemyBigCollisionRadius = 26
	enemyBigHitCooldown     = time.Millisecond * 250
	enemyBigShotCooldown    = time.Millisecond * 2500
	enemyBigMaxHitCount     = 4
	enemyBigBulletSpeed     = 5
)

type enemyBig struct {
	renderer        *sdl.Renderer
	tex             *sdl.Texture
	texHit          *sdl.Texture
	x, y            float64
	active          bool
	hit             bool
	lastTimeHit     time.Time
	lastTimeBooster time.Time
	lastTimeShot    time.Time
	texXPos         int32
	hitCount        int8
}

func newEnemyBig(renderer *sdl.Renderer, x, y float64) *enemyBig {

	en := enemyBig{
		renderer: renderer,
		tex:      newTexture(renderer, "sprites/enemy-big.png"),
		texHit:   newTexture(renderer, "sprites/enemy-big-hit.png"),
		texXPos:  0,
		x:        x,
		y:        y,
		active:   false}

	return &en
}

func (e *enemyBig) start(x, y, angle, speed float64, entityType int8) {

	// TODO receive values
	e.x = float64(rand.Intn(screenWidth))
	e.y = -30
	e.hitCount = 0
	e.active = true
}

func (e *enemyBig) executeCollisionWith(other entity) {

	if other.getType() == entityTypePlayerBullet {
		e.beHit()
	} else if other.getType() == entityTypePlayer {
		e.beDestroyed()
	}
}

func (e *enemyBig) getCollisionCircle() circle {

	return circle{x: e.x, y: e.y, radius: enemyBigCollisionRadius}
}

func (e *enemyBig) isActive() bool {

	return e.active
}

func (e *enemyBig) getType() int8 {

	return entityTypeEnemyBig
}

func (e *enemyBig) draw() {

	if !e.active {
		return
	}

	drawEnemySize := enemyBigSize * scale
	drawX := e.x - float64(drawEnemySize)/2.0
	drawY := e.y - float64(drawEnemySize)/2.0
	drawRect := sdl.Rect{X: int32(drawX), Y: int32(drawY), W: int32(drawEnemySize), H: int32(drawEnemySize)}

	var tex *sdl.Texture
	if e.hit {
		tex = e.texHit
	} else {
		tex = e.tex
	}

	xTex := e.texXPos * enemyBigSize

	e.renderer.CopyEx(tex,
		&sdl.Rect{X: xTex, Y: 0, W: enemyBigSize, H: enemyBigSize},
		&drawRect,
		0,
		&sdl.Point{X: drawRect.W / 2.0, Y: drawRect.H / 2.0},
		sdl.FLIP_NONE)

	if deb.active {
		// debug drawRect
		e.renderer.SetDrawColor(255, 0, 0, 255)
		e.renderer.DrawRect(&drawRect)

		// debug rect of collision
		collisionRect := sdl.Rect{X: int32(e.x - enemyBigCollisionRadius), Y: int32(e.y - enemyBigCollisionRadius),
			W: enemyBigCollisionRadius * 2, H: enemyBigCollisionRadius * 2}
		e.renderer.DrawRect(&collisionRect)
	}
}

func (e *enemyBig) update() {

	if !e.active {
		return
	}

	if e.hit && time.Since(e.lastTimeHit) >= enemyBigHitCooldown {
		e.hit = false
	}

	if time.Since(e.lastTimeBooster) >= playerShotCooldown/2 {
		if e.texXPos == 0 {
			e.texXPos = 1
		} else if e.texXPos == 1 {
			e.texXPos = 0
		}
		e.lastTimeBooster = time.Now()
	}

	e.shoot()

	eSpeed := enemyBigSpeed * delta

	e.y += eSpeed
	if e.y > screenHeight {
		e.active = false
	}
}

func (e *enemyBig) beHit() {

	score.incrementPointsP1(1)
	chunkHit.Play(0, 0)
	e.hit = true
	e.hitCount++
	e.lastTimeHit = time.Now()

	if e.hitCount >= enemyBigMaxHitCount {
		e.beDestroyed()
	}
}

func (e *enemyBig) beDestroyed() {

	chunkExplosion.Play(2, 0)
	e.active = false
	ex := explosionFromPool()
	ex.start(e.x, e.y, enemyBigSpeed)
}

func (e *enemyBig) shoot() {

	if time.Since(e.lastTimeShot) >= enemyBigShotCooldown {
		chunkLaser.Play(1, 0)
		bul := bulletFromPool()
		if bul != nil {
			p, ok := plr.(*player)
			if ok {
				angle := angleOfLine(e.x, e.y, p.x, p.y)
				bul.start(e.x, e.y, angle, enemyBigBulletSpeed, entityTypeEnemyBullet)
			}
			e.lastTimeShot = time.Now()
		}
	}
}

func (e *enemyBig) deactivate() {

	e.active = false
}

var enemiesBig []entity

func initEnemiesBig(renderer *sdl.Renderer) {
	for i := 0; i < 20; i++ {
		en := newEnemyBig(renderer, screenWidth/2+enemyBigSize, -1*enemyBigSize)
		enemiesBig = append(enemiesBig, en)
	}
}

func deactivateAllEnemiesBig() {
	for _, en := range enemiesBig {
		en.deactivate()
	}
}

func enemyBigFromPool() entity {
	for _, en := range enemiesBig {
		if !en.isActive() {
			return en
		}
	}

	return nil
}
