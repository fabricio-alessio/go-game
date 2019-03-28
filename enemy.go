package main

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	enemySpeed           = 1
	enemySize            = 32
	enemyCollisionRadius = 26
	enemyHitCooldown     = time.Millisecond * 250
	enemyShotCooldown    = time.Millisecond * 2500
	enemyMaxHitCount     = 4
	enemyBulletSpeed     = 5
)

type enemy struct {
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

func newEnemy(renderer *sdl.Renderer, x, y float64) (e enemy) {

	e.tex = newTexture(renderer, "sprites/enemy-big.png")
	e.texHit = newTexture(renderer, "sprites/enemy-big-hit.png")
	e.x = x
	e.y = y
	e.active = false

	e.texXPos = 0

	return e
}

func (e *enemy) draw(renderer *sdl.Renderer) {

	if !e.active {
		return
	}

	drawEnemySize := enemySize * scale
	drawX := e.x - float64(drawEnemySize)/2.0
	drawY := e.y - float64(drawEnemySize)/2.0
	drawRect := sdl.Rect{X: int32(drawX), Y: int32(drawY), W: int32(drawEnemySize), H: int32(drawEnemySize)}

	var tex *sdl.Texture
	if e.hit {
		tex = e.texHit
	} else {
		tex = e.tex
	}

	xTex := e.texXPos * enemySize

	renderer.CopyEx(tex,
		&sdl.Rect{X: xTex, Y: 0, W: enemySize, H: enemySize},
		&drawRect,
		0,
		&sdl.Point{X: drawRect.W / 2.0, Y: drawRect.H / 2.0},
		sdl.FLIP_NONE)

	if deb.active {
		// debug drawRect
		renderer.SetDrawColor(255, 0, 0, 255)
		renderer.DrawRect(&drawRect)

		// debug rect of collision
		collisionRect := sdl.Rect{X: int32(e.x - enemyCollisionRadius), Y: int32(e.y - enemyCollisionRadius),
			W: enemyCollisionRadius * 2, H: enemyCollisionRadius * 2}
		renderer.DrawRect(&collisionRect)
	}
}

func (e *enemy) update() {

	if !e.active {
		return
	}

	if e.hit && time.Since(e.lastTimeHit) >= enemyHitCooldown {
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

	eSpeed := enemySpeed * delta

	e.y += eSpeed
	if e.y > screenHeight {
		e.active = false
	}
}

func (e *enemy) beHit() {

	score.increment(1)
	chunkHit.Play(0, 0)
	e.hit = true
	e.hitCount++
	e.lastTimeHit = time.Now()

	if e.hitCount >= enemyMaxHitCount {
		e.beDestroyed()
	}
}

func (e *enemy) beDestroyed() {

	chunkExplosion.Play(2, 0)
	e.active = false
	ex := explosionFromPool()
	ex.start(e.x, e.y, enemySpeed)
}

func (e *enemy) shoot() {
	if time.Since(e.lastTimeShot) >= enemyShotCooldown {
		chunkLaser.Play(1, 0)
		bul := bulletFromPool()
		if bul != nil {
			angle := angleOfLine(e.x, e.y, plr.x, plr.y)
			bul.start(e.x, e.y, angle, enemyBulletSpeed, false)
			e.lastTimeShot = time.Now()
		}
	}
}

var enemies []*enemy

func initEnemies(renderer *sdl.Renderer) {
	for i := 0; i < 20; i++ {
		en := newEnemy(renderer, screenWidth/2+enemySize, -1*enemySize)
		enemies = append(enemies, &en)
	}
}

func deactivateAllEnemies() {
	for _, en := range enemies {
		en.active = false
	}
}
