package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	explosionSize      = 16
	explosionCooldown  = time.Millisecond * 100
	explosionFaseCount = 5
)

type explosion struct {
	tex             *sdl.Texture
	x, y            float64
	active          bool
	lastTimeBooster time.Time
	texXPos         int32
	speed           float64
	angle           float64
}

func newExplosion(renderer *sdl.Renderer) (ex explosion) {

	ex.tex = newTexture(renderer, "sprites/explosion.png")
	ex.active = false

	ex.texXPos = 0
	return ex
}

func (ex *explosion) draw(renderer *sdl.Renderer) {

	if !ex.active {
		return
	}

	w := int32(explosionSize * scale)
	h := int32(explosionSize * scale)
	x := ex.x - (explosionSize*scale)/2.0
	y := ex.y - (explosionSize*scale)/2.0

	xTex := ex.texXPos * explosionSize

	renderer.CopyEx(ex.tex,
		&sdl.Rect{X: xTex, Y: 0, W: explosionSize, H: explosionSize},
		&sdl.Rect{X: int32(x), Y: int32(y), W: w, H: h},
		ex.angle,
		&sdl.Point{X: w / 2.0, Y: h / 2.0},
		sdl.FLIP_NONE)
}

func (ex *explosion) update() {

	if !ex.active {
		return
	}

	if time.Since(ex.lastTimeBooster) >= explosionCooldown {
		if ex.texXPos >= explosionFaseCount {
			ex.texXPos = 0
			ex.active = false
		} else {
			ex.texXPos++
		}

		ex.lastTimeBooster = time.Now()
	}

	exSpeed := ex.speed * delta

	ex.y += exSpeed
}

func (ex *explosion) start(x, y, speed float64) {

	ex.active = true
	ex.x = x
	ex.y = y
	ex.speed = speed
	ex.texXPos = 0
	ex.angle = float64(rand.Intn(4) * 90)
	deb.set(0, fmt.Sprintf("angle %f", ex.angle))
	ex.lastTimeBooster = time.Now()
}

var explosions []*explosion

func initExplosions(renderer *sdl.Renderer) {
	for i := 0; i < 20; i++ {
		ex := newExplosion(renderer)
		explosions = append(explosions, &ex)
	}
}

func explosionFromPool() *explosion {
	for _, ex := range explosions {
		if !ex.active {
			return ex
		}
	}

	return nil
}
