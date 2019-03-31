package main

import (
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

func newenemyBig(renderer *sdl.Renderer, x, y float64) (e enemyBig) {

	e.tex = newTexture(renderer, "sprites/enemy-big.png")
	e.texHit = newTexture(renderer, "sprites/enemy-big-hit.png")
	e.x = x
	e.y = y
	e.active = false

	e.texXPos = 0

	return e
}

func (e *enemyBig) draw(renderer *sdl.Renderer) {

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

	renderer.CopyEx(tex,
		&sdl.Rect{X: xTex, Y: 0, W: enemyBigSize, H: enemyBigSize},
		&drawRect,
		0,
		&sdl.Point{X: drawRect.W / 2.0, Y: drawRect.H / 2.0},
		sdl.FLIP_NONE)

	if deb.active {
		// debug drawRect
		renderer.SetDrawColor(255, 0, 0, 255)
		renderer.DrawRect(&drawRect)

		// debug rect of collision
		collisionRect := sdl.Rect{X: int32(e.x - enemyBigCollisionRadius), Y: int32(e.y - enemyBigCollisionRadius),
			W: enemyBigCollisionRadius * 2, H: enemyBigCollisionRadius * 2}
		renderer.DrawRect(&collisionRect)
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
			angle := angleOfLine(e.x, e.y, plr.x, plr.y)
			bul.start(e.x, e.y, angle, enemyBigBulletSpeed, false)
			e.lastTimeShot = time.Now()
		}
	}
}

var enemiesBig []*enemyBig

func initEnemiesBig(renderer *sdl.Renderer) {
	for i := 0; i < 20; i++ {
		en := newenemyBig(renderer, screenWidth/2+enemyBigSize, -1*enemyBigSize)
		enemiesBig = append(enemiesBig, &en)
	}
}

func deactivateAllEnemiesBig() {
	for _, en := range enemiesBig {
		en.active = false
	}
}

func enemyBigFromPool() *enemyBig {
	for _, en := range enemiesBig {
		if !en.active {
			return en
		}
	}

	return nil
}
