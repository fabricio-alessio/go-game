package main

import (
	"fmt"
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	playerSpeed           = 5
	playerWidth           = 16
	playerHeight          = 24
	playerShotCooldown    = time.Millisecond * 250
	playerCollisionRadius = 16
	playerBulletSpeed     = 10
	playerInitialLives    = 3
	playerMaxShield       = 5
)

type player struct {
	renderer        *sdl.Renderer
	tex             *sdl.Texture
	x, y            float64
	lastTimeShot    time.Time
	lastTimeDrive   time.Time
	lastTimeUndrive time.Time
	lastTimeBooster time.Time
	texXPos         int32
	texYPos         int32
	lives           int8
	shield          int8
	active          bool
	starting        bool
}

func newPlayer(renderer *sdl.Renderer) *player {

	plr := player{
		renderer: renderer,
		tex:      newTexture(renderer, "sprites/ship.png"),
		x:        screenWidth / 2.0,
		y:        screenHeight - (playerHeight * scale),
		texXPos:  2,
		texYPos:  1,
		active:   true}

	plr.setLives(playerInitialLives)
	plr.setShield(playerMaxShield)

	return &plr
}

func (p *player) draw() {

	if !p.active {
		return
	}

	drawPlayerWidth := playerWidth * scale
	drawPlayerHeight := playerHeight * scale
	drawX := p.x - float64(drawPlayerWidth)/2.0
	drawY := p.y - float64(drawPlayerHeight)/2.0
	drawRect := sdl.Rect{X: int32(drawX), Y: int32(drawY), W: int32(drawPlayerWidth), H: int32(drawPlayerHeight)}

	xTex := p.texXPos * playerWidth
	yTex := p.texYPos * playerHeight

	p.renderer.Copy(p.tex,
		&sdl.Rect{X: xTex, Y: yTex, W: playerWidth, H: playerHeight},
		&drawRect)

	if deb.active {
		// debug drawRect
		p.renderer.SetDrawColor(255, 0, 0, 255)
		p.renderer.DrawRect(&drawRect)

		// debug rect of collision
		collisionRect := sdl.Rect{X: int32(p.x - playerCollisionRadius), Y: int32(p.y - playerCollisionRadius),
			W: playerCollisionRadius * 2, H: playerCollisionRadius * 2}
		p.renderer.DrawRect(&collisionRect)
	}
}

func (p *player) moveUp() {
	if p.y-(playerHeight/2.0) > 0 {
		p.y -= playerSpeed * delta
	}
}

func (p *player) moveDown() {
	if p.y+(playerHeight/2.0) < screenHeight {
		p.y += playerSpeed * delta
	}
}

func (p *player) moveLeft() {
	if p.x-(playerWidth/2.0) > 0 {
		p.x -= playerSpeed * delta
		if time.Since(p.lastTimeDrive) >= playerShotCooldown {
			if p.texXPos != 0 {
				p.texXPos -= 1
			}
			p.lastTimeDrive = time.Now()
		}
	}
}

func (p *player) moveRight() {
	if p.x+(playerWidth/2.0) < screenWidth {
		p.x += playerSpeed * delta
		if time.Since(p.lastTimeDrive) >= playerShotCooldown {
			if p.texXPos != 4 {
				p.texXPos += 1
			}
			p.lastTimeDrive = time.Now()
		}
	}
}

func (p *player) moveCenter() {
	p.texXPos = 2
	p.lastTimeUndrive = time.Now()
}

func (p *player) shoot() {
	if time.Since(p.lastTimeShot) >= playerShotCooldown {
		chunkLaser.Play(1, 0)
		bul := bulletFromPool()
		if bul != nil {
			bul.start(p.x, p.y, 270*(math.Pi/180), playerBulletSpeed, entityTypePlayerBullet)
			p.lastTimeShot = time.Now()
		}
	}
}

func (p *player) handleKeyboard() {

	keys := sdl.GetKeyboardState()

	if keys[sdl.SCANCODE_LEFT] == 1 {
		p.moveLeft()
	} else if keys[sdl.SCANCODE_RIGHT] == 1 {
		p.moveRight()
	} else if keys[sdl.SCANCODE_LEFT] == 0 && keys[sdl.SCANCODE_RIGHT] == 0 {
		p.moveCenter()
	}

	if keys[sdl.SCANCODE_UP] == 1 {
		p.moveUp()
	} else if keys[sdl.SCANCODE_DOWN] == 1 {
		p.moveDown()
	}

	if keys[sdl.SCANCODE_SPACE] == 1 {
		p.shoot()
	}
}

func (p *player) handleJoystic() {

	if joy.Button(2) > 0 {
		p.shoot()
	}
	if joy.Axis(0) < -8000 {
		p.moveLeft()
	} else if joy.Axis(0) > 8000 {
		p.moveRight()
	}
	if joy.Axis(1) < -8000 {
		p.moveUp()
	} else if joy.Axis(1) > 8000 {
		p.moveDown()
	}
}

func (p *player) update() {

	if !p.active {
		return
	}

	if p.starting {
		p.moveUp()
		if p.y <= screenHeight-playerHeight*4 {
			p.starting = false
		}
	} else {

		if time.Since(p.lastTimeBooster) >= playerShotCooldown/2 {
			if p.texYPos == 0 {
				p.texYPos = 1
			} else if p.texYPos == 1 {
				p.texYPos = 0
			}
			p.lastTimeBooster = time.Now()
		}

		p.handleKeyboard()
		p.handleJoystic()
	}
}

func (p *player) beHit() {

	p.setShield(p.shield - 1)
	chunkHit.Play(0, 0)

	if p.shield <= 0 {
		p.beDestroyed()
	}
	deb.set(1, fmt.Sprintf("shield %d lives %d", p.shield, p.lives))
}

func (p *player) beDestroyed() {

	chunkExplosion.Play(2, 0)
	p.setLives(p.lives - 1)
	ex := explosionFromPool()
	ex.start(p.x, p.y, enemyBigSpeed)

	p.active = false
	if p.lives > 0 {
		efleet.askRestartLevel()
	} else {
		efleet.askGoToMenu()
	}
}

func (p *player) replay() {

	p.setLives(playerInitialLives)
}

func (p *player) setLives(qtd int8) {

	p.lives = qtd
	score.setLivesP1(qtd)
}

func (p *player) setShield(qtd int8) {

	p.shield = qtd
	score.setShieldP1(qtd)
}

func (p *player) start(x, y, angle, speed float64, entityType int8) {

	p.x = screenWidth / 2.0
	p.y = screenHeight + (playerHeight * scale)
	p.setShield(playerMaxShield)
	p.active = true
	p.texXPos = 2
	p.texYPos = 1
	p.starting = true
}

func (p *player) executeCollisionWith(other entity) {

	if other.getType() == entityTypeEnemyBullet {
		p.beHit()
	} else if other.getType() == entityTypeEnemyBig {
		p.beDestroyed()
	} else if other.getType() == entityTypeEnemySmall {
		p.beDestroyed()
	}
}

func (p *player) getCollisionCircle() circle {

	return circle{x: p.x, y: p.y, radius: playerCollisionRadius}
}

func (p *player) isActive() bool {

	return p.active
}

func (p *player) getType() int8 {

	return entityTypePlayer
}

func (p *player) deactivate() {

	p.active = false
}
