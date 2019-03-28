package main

import (
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	bulletSize            = 16
	bulletCollisionRadius = 6
)

type bullet struct {
	tex             *sdl.Texture
	x, y            float64
	angle           float64
	speed           float64
	active          bool
	texXPos         int32
	texYPos         int32
	lastTimeBooster time.Time
	fromPlayer      bool
}

func newBullet(renderer *sdl.Renderer) (b bullet) {

	b.tex = newTexture(renderer, "sprites/laser-bolts2.png")

	b.texXPos = 0
	b.texYPos = 1

	return b
}

func (b *bullet) draw(renderer *sdl.Renderer) {

	if !b.active {
		return
	}

	drawBulletSize := bulletSize * scale
	drawX := b.x - float64(drawBulletSize)/2.0
	drawY := b.y - float64(drawBulletSize)/2.0
	drawRect := sdl.Rect{X: int32(drawX), Y: int32(drawY), W: int32(drawBulletSize), H: int32(drawBulletSize)}

	xTex := b.texXPos * bulletSize
	yTex := b.texYPos * bulletSize

	renderer.Copy(b.tex,
		&sdl.Rect{X: xTex, Y: yTex, W: bulletSize, H: bulletSize},
		&drawRect)

	if deb.active {
		// debug drawRect
		renderer.SetDrawColor(255, 0, 0, 255)
		renderer.DrawRect(&drawRect)

		// debug rect of collision
		collisionRect := sdl.Rect{X: int32(b.x - bulletCollisionRadius), Y: int32(b.y - bulletCollisionRadius),
			W: bulletCollisionRadius * 2, H: bulletCollisionRadius * 2}
		renderer.DrawRect(&collisionRect)
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

func (b *bullet) start(x, y, angle, speed float64, fromPlayer bool) {

	b.x = x
	b.y = y
	b.angle = angle
	b.speed = speed
	b.fromPlayer = fromPlayer
	if fromPlayer {
		b.texYPos = 1
	} else {
		b.texYPos = 0
	}

	b.active = true
}

var bulletPool []*bullet

func initBulletPool(renderer *sdl.Renderer) {

	for i := 0; i < 30; i++ {
		b := newBullet(renderer)
		bulletPool = append(bulletPool, &b)
	}
}

func bulletFromPool() *bullet {
	for _, bul := range bulletPool {
		if !bul.active {
			return bul
		}
	}

	return nil
}

func deactivateAllBullets() {
	for _, bul := range bulletPool {
		bul.active = false
	}
}
