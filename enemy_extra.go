package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	enemyExtraSpeed           = 1
	enemyExtraSize            = 46
	enemyExtraCollisionRadius = 36
	enemyExtraHitCooldown     = time.Millisecond * 250
	enemyExtraShotCooldown    = time.Millisecond * 2500
	enemyExtraMaxHitCount     = 7
	enemyExtraBulletSpeed     = 5
)

type enemyExtra struct {
	renderer        *sdl.Renderer
	tex             *sdl.Texture
	x, y            float64
	active, hit     bool
	lastTimeHit     time.Time
	lastTimeBooster time.Time
	lastTimeShot    time.Time
	texXPos         int32
	hitCount        int8
}

func newEnemyExtra(renderer *sdl.Renderer, x, y float64) *enemyExtra {

	en := enemyExtra{
		renderer: renderer,
		tex:      newTexture(renderer, "sprites/enemy-extra.png"),
		texXPos:  0,
		x:        x,
		y:        y,
		active:   false}

	return &en
}

func (e *enemyExtra) start(x, y, angle, speed float64, entityType int8) {

	e.x = x
	e.y = -30
	e.hitCount = 0
	e.active = true
}

func (e *enemyExtra) executeCollisionWith(other entity) {

	if other.getType() == entityTypePlayerBullet {
		e.beHit()
	} else if other.getType() == entityTypePlayer {
		e.beDestroyed()
	}
}

func (e *enemyExtra) getCollisionCircle() circle {

	return circle{x: e.x, y: e.y, radius: enemyExtraCollisionRadius}
}

func (e *enemyExtra) isActive() bool {

	return e.active
}

func (e *enemyExtra) getType() int8 {

	return entityTypeEnemyExtra
}

func (e *enemyExtra) draw() {

	if !e.active {
		return
	}

	drawEnemySize := enemyExtraSize * scale
	drawX := e.x - float64(drawEnemySize)/2.0
	drawY := e.y - float64(drawEnemySize)/2.0
	drawRect := sdl.Rect{X: int32(drawX), Y: int32(drawY), W: int32(drawEnemySize), H: int32(drawEnemySize)}

	xTex := e.texXPos * enemyExtraSize

	e.renderer.CopyEx(e.tex,
		&sdl.Rect{X: xTex, Y: 0, W: enemyExtraSize, H: enemyExtraSize},
		&drawRect,
		0,
		&sdl.Point{X: drawRect.W / 2.0, Y: drawRect.H / 2.0},
		sdl.FLIP_NONE)

	if deb.active {
		// debug drawRect
		e.renderer.SetDrawColor(255, 0, 0, 255)
		e.renderer.DrawRect(&drawRect)

		// debug rect of collision
		collisionRect := sdl.Rect{X: int32(e.x - float64(enemyExtraCollisionRadius)), Y: int32(e.y - float64(enemyExtraCollisionRadius)),
			W: enemyExtraCollisionRadius * 2, H: enemyExtraCollisionRadius * 2}
		e.renderer.DrawRect(&collisionRect)
	}
}

func (e *enemyExtra) booster() {

	if time.Since(e.lastTimeBooster) >= playerShotCooldown/2 {
		if e.texXPos == 0 {
			e.texXPos = 1
		} else if e.texXPos == 1 {
			e.texXPos = 0
		}
		e.lastTimeBooster = time.Now()
	}
}

func (e *enemyExtra) update() {

	if !e.active {
		return
	}

	if e.hit && time.Since(e.lastTimeHit) >= enemyBigHitCooldown {
		e.hit = false
	}

	e.booster()
	e.shoot()

	e.y += enemyExtraSpeed * delta
	if e.y > screenHeight {
		fmt.Println("Send extra enemy back to poll")
		e.active = false
	}
}

func (e *enemyExtra) beHit() {

	score.incrementPointsP1(1)
	mixer.playSound("hit")
	e.hit = true
	e.hitCount++
	e.lastTimeHit = time.Now()

	if e.hitCount >= enemyExtraMaxHitCount {
		e.beDestroyed()
	}
}

func (e *enemyExtra) beDestroyed() {

	mixer.playSound("explosion")
	e.active = false
	ex := explosionFromPool()
	ex.start(e.x, e.y, enemyExtraSpeed)
	pu := powerUpFromPool()
	pu.start(e.x, 0, 0, 0, entityTypePowerUpHelper)
}

func (e *enemyExtra) shoot() {

	if time.Since(e.lastTimeShot) >= enemyBigShotCooldown {
		mixer.playSound("laser")
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

func (e *enemyExtra) deactivate() {

	e.active = false
}

var enemiesExtra []entity

func initEnemiesExtra(renderer *sdl.Renderer) {
	for i := 0; i < 20; i++ {
		en := newEnemyExtra(renderer, screenWidth/2+enemyExtraSize, -1*enemyExtraSize)
		enemiesExtra = append(enemiesExtra, en)
	}
}

func deactivateAllEnemiesExtra() {
	for _, en := range enemiesExtra {
		en.deactivate()
	}
}

func enemyExtraFromPool() entity {
	for _, en := range enemiesExtra {
		if !en.isActive() {
			return en
		}
	}

	return nil
}
