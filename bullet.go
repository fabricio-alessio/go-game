package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	bulletSize            = 16
	bulletCollisionRadius = 6
)

type bullet struct {
	renderer        *sdl.Renderer
	tex             *sdl.Texture
	x, y            float64
	angle           float64
	speed           float64
	active          bool
	texXPos         int32
	texYPos         int32
	lastTimeBooster time.Time
	entityType      int8
}

func newBullet(renderer *sdl.Renderer) *bullet {

	bul := bullet{
		renderer: renderer,
		tex:      newTexture(renderer, "sprites/laser-bolts2.png"),
		texXPos:  0,
		texYPos:  1}

	return &bul
}

func (b *bullet) deactivate() {

	b.active = false
}

func (b *bullet) getCollisionCircle() circle {

	return circle{x: b.x, y: b.y, radius: bulletCollisionRadius}
}

func (b *bullet) isActive() bool {

	return b.active
}

func (b *bullet) getType() int8 {

	return b.entityType
}

func (b *bullet) executeCollisionWith(other entity) {

	b.deactivate()

	for i := 0; i < 3; i++ {
		b.startBlast()
	}
}

func (b *bullet) startBlast() {

	speedVar := b.speed/2 + rand.Float64()*b.speed
	blast := bulletBlastFromPool()
	degreeVar := rand.Intn(90) - 45
	radVar := float64(degreeVar) * math.Pi / 180
	blast.start(b.x, b.y, b.angle+radVar, speedVar, 0)
}

func (b *bullet) draw() {

	if !b.active {
		return
	}

	drawBulletSize := bulletSize * scale
	drawX := b.x - float64(drawBulletSize)/2.0
	drawY := b.y - float64(drawBulletSize)/2.0
	drawRect := sdl.Rect{X: int32(drawX), Y: int32(drawY), W: int32(drawBulletSize), H: int32(drawBulletSize)}

	xTex := b.texXPos * bulletSize
	yTex := b.texYPos * bulletSize

	b.renderer.Copy(b.tex,
		&sdl.Rect{X: xTex, Y: yTex, W: bulletSize, H: bulletSize},
		&drawRect)

	if deb.active {
		// debug drawRect
		b.renderer.SetDrawColor(255, 0, 0, 255)
		b.renderer.DrawRect(&drawRect)

		// debug rect of collision
		collisionRect := sdl.Rect{X: int32(b.x - bulletCollisionRadius), Y: int32(b.y - bulletCollisionRadius),
			W: bulletCollisionRadius * 2, H: bulletCollisionRadius * 2}
		b.renderer.DrawRect(&collisionRect)
	}
}

func (b *bullet) update() {

	if !b.active {
		return
	}

	if time.Since(b.lastTimeBooster) >= playerShotCooldown/2 {
		if b.texXPos == 0 {
			b.texXPos = 1
		} else if b.texXPos == 1 {
			b.texXPos = 0
		}
		b.lastTimeBooster = time.Now()
	}

	bSpeed := b.speed * delta
	b.x += bSpeed * math.Cos(b.angle)
	b.y += bSpeed * math.Sin(b.angle)

	if b.x <= -bulletSize || b.y <= -bulletSize || b.x >= screenWidth+bulletSize || b.y >= screenHeight+bulletSize {
		b.active = false
	}
}

func (b *bullet) start(x, y, angle, speed float64, entityType int8) {

	b.x = x
	b.y = y
	b.angle = angle
	b.speed = speed
	b.entityType = entityType
	if entityType == entityTypePlayerBullet {
		b.texYPos = 1
	} else {
		b.texYPos = 0
	}

	b.active = true
}

var bulletPool []entity

func initBulletPool(renderer *sdl.Renderer) {

	for i := 0; i < 60; i++ {
		b := newBullet(renderer)
		bulletPool = append(bulletPool, b)
	}
}

func bulletFromPool() entity {
	for _, bul := range bulletPool {
		if !bul.isActive() {
			return bul
		}
	}

	return nil
}

func deactivateAllBullets() {
	for _, bul := range bulletPool {
		bul.deactivate()
	}
}
