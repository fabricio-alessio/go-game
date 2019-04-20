package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	bulletBlastSize = 2
)

type bulletBlast struct {
	renderer *sdl.Renderer
	x, y     float64
	angle    float64
	speed    float64
	active   bool
	power    float64
}

func newBulletBlast(renderer *sdl.Renderer) *bulletBlast {

	bul := bulletBlast{
		renderer: renderer,
		power:    255}

	return &bul
}

func (b *bulletBlast) deactivate() {

	b.active = false
}

func (b *bulletBlast) getCollisionCircle() circle {

	return circle{}
}

func (b *bulletBlast) isActive() bool {

	return b.active
}

func (b *bulletBlast) getType() int8 {

	return entityTypeBulletBlast
}

func (b *bulletBlast) executeCollisionWith(other entity) {

	if b.getType() == entityTypeEnemyBullet && other.getType() == entityTypePlayer {
		b.deactivate()
	} else if b.getType() == entityTypePlayerBullet && other.getType() == entityTypeEnemyBig {
		b.deactivate()
	} else if b.getType() == entityTypePlayerBullet && other.getType() == entityTypeEnemySmall {
		b.deactivate()
	} else if b.getType() == entityTypePlayerBullet && other.getType() == entityTypeEnemyExtra {
		b.deactivate()
	}
}

func (b *bulletBlast) draw() {

	if !b.active {
		return
	}

	drawBulletSize := bulletBlastSize * scale
	drawX := b.x - float64(drawBulletSize)/2.0
	drawY := b.y - float64(drawBulletSize)/2.0
	drawRect := sdl.Rect{X: int32(drawX), Y: int32(drawY), W: int32(drawBulletSize), H: int32(drawBulletSize)}

	b.renderer.SetDrawColor(uint8(b.power), uint8(b.power), uint8(b.power), uint8(b.power))
	b.renderer.FillRect(&drawRect)
}

func (b *bulletBlast) update() {

	if !b.active {
		return
	}

	b.power = b.power - 10*delta

	bSpeed := b.speed * delta
	b.x += bSpeed * math.Cos(b.angle)
	b.y += bSpeed * math.Sin(b.angle)

	if b.x <= -bulletSize || b.y <= -bulletSize || b.x >= screenWidth+bulletSize || b.y >= screenHeight+bulletSize {
		b.active = false
	}
	if b.power <= 0 {
		b.active = false
	}
}

func (b *bulletBlast) start(x, y, angle, speed float64, entityType int8) {

	b.x = x
	b.y = y
	b.angle = angle
	b.speed = speed
	b.power = 255

	b.active = true
}

var bulletBlastPool []entity

func initBulletBlastPool(renderer *sdl.Renderer) {

	for i := 0; i < 80; i++ {
		b := newBulletBlast(renderer)
		bulletBlastPool = append(bulletBlastPool, b)
	}
}

func bulletBlastFromPool() entity {
	for _, bul := range bulletBlastPool {
		if !bul.isActive() {
			return bul
		}
	}

	return nil
}

func deactivateAllBulletBlasts() {
	for _, bul := range bulletBlastPool {
		bul.deactivate()
	}
}
