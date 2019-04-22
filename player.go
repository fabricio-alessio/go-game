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
	playerHelperSize      = 19
	playerShotCooldown    = time.Millisecond * 250
	playerCollisionRadius = 16
	playerBulletSpeed     = 10
	playerInitialLives    = 3
	playerMaxShield       = 5
	playerMovSize         = 500
	playerMaxHelper       = 3
	playerMaxPower        = 6
)

type position struct {
	x, y, angle float64
}

type player struct {
	renderer              *sdl.Renderer
	tex                   *sdl.Texture
	x, y                  float64
	lastTimeShot          time.Time
	lastTimeDrive         time.Time
	lastTimeUndrive       time.Time
	lastTimeBooster       time.Time
	lastTimeHelperBooster time.Time
	lastTimeRegisterMov   time.Time
	texXPos               int32
	texYPos               int32
	lives                 int8
	shield                int8
	active                bool
	starting              bool
	power                 int8
	moviment              [playerMovSize]position
	movIndex              int
	texHelper             *sdl.Texture
	helperX, helperY      [playerMaxHelper]float64
	helperCount           int
	helperTexXPos         int32
}

func newPlayer(renderer *sdl.Renderer) *player {

	plr := player{
		renderer:  renderer,
		tex:       newTexture(renderer, "sprites/ship.png"),
		texHelper: newTexture(renderer, "sprites/secundary.png"),
		x:         screenWidth / 2.0,
		y:         screenHeight - (playerHeight * scale),
		texXPos:   2,
		texYPos:   1,
		active:    true,
		power:     0}

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

	for i := 0; i < p.helperCount; i++ {
		p.drawHelper(i)
	}
}

func (p *player) drawHelper(index int) {

	helperSize := playerHelperSize * scale
	drawX := p.helperX[index] - float64(helperSize)/2.0
	drawY := p.helperY[index] - float64(helperSize)/2.0
	drawRect := sdl.Rect{X: int32(drawX), Y: int32(drawY), W: int32(helperSize), H: int32(helperSize)}

	xTex := p.helperTexXPos * playerHelperSize

	p.renderer.Copy(p.texHelper,
		&sdl.Rect{X: xTex, Y: 0, W: playerHelperSize, H: playerHelperSize},
		&drawRect)
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
		mixer.playSound("laser")
		if p.power == 0 {
			startBullet(p.x, p.y, 270*(math.Pi/180))
		} else if p.power == 1 {
			startBullet(p.x-10, p.y, 270*(math.Pi/180))
			startBullet(p.x+10, p.y, 270*(math.Pi/180))
		} else if p.power == 2 {
			startBullet(p.x-15, p.y, 255*(math.Pi/180))
			startBullet(p.x, p.y, 270*(math.Pi/180))
			startBullet(p.x+15, p.y, 285*(math.Pi/180))
		} else if p.power >= 3 {
			startBullet(p.x-20, p.y, 240*(math.Pi/180))
			startBullet(p.x-15, p.y, 255*(math.Pi/180))
			startBullet(p.x, p.y, 270*(math.Pi/180))
			startBullet(p.x+15, p.y, 285*(math.Pi/180))
			startBullet(p.x+20, p.y, 300*(math.Pi/180))
		}
		doubleCount := p.power - 3
		for i := 0; i < p.helperCount; i++ {
			if doubleCount <= 0 {
				startBullet(p.helperX[i], p.helperY[i], 270*(math.Pi/180))
			} else {
				startBullet(p.helperX[i]-10, p.helperY[i], 270*(math.Pi/180))
				startBullet(p.helperX[i]+10, p.helperY[i], 270*(math.Pi/180))
			}
			doubleCount--
		}
		p.lastTimeShot = time.Now()
	}
}

func startBullet(x, y, angle float64) {
	bul := bulletFromPool()
	if bul != nil {
		bul.start(x, y, angle, playerBulletSpeed, entityTypePlayerBullet)
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

	if joystic0.Button(2) > 0 {
		p.shoot()
	}
	if joystic0.Axis(0) < -8000 {
		p.moveLeft()
	} else if joystic0.Axis(0) > 8000 {
		p.moveRight()
	}
	if joystic0.Axis(1) < -8000 {
		p.moveUp()
	} else if joystic0.Axis(1) > 8000 {
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

	p.registerMoviment()
	for i := 0; i < p.helperCount; i++ {
		p.updateHelper(i)
	}
}

func (p *player) registerMoviment() {

	//if time.Since(p.lastTimeRegisterMov) >= playerShotCooldown/8 {

	pos := position{x: p.x, y: p.y, angle: 0}
	if p.moviment[p.movIndex] != pos {
		if p.movIndex == playerMovSize-1 {
			p.movIndex = 0
		} else {
			p.movIndex++
		}
		p.moviment[p.movIndex] = pos
	}

	//	p.lastTimeRegisterMov = time.Now()
	//}
}

func (p *player) getPositionDelayed(qtd int) position {

	index := p.movIndex - qtd
	if index < 0 {
		index = playerMovSize - 1 + index
	}
	return p.moviment[index]
}

func (p *player) updateHelper(index int) {

	pos := p.getPositionDelayed(10 + index*10)
	p.helperX[index] = pos.x
	p.helperY[index] = pos.y

	if time.Since(p.lastTimeHelperBooster) >= playerShotCooldown/4 {
		if p.helperTexXPos == 7 {
			p.helperTexXPos = 0
		} else {
			p.helperTexXPos++
		}
		p.lastTimeHelperBooster = time.Now()
	}
}

func (p *player) beHit() {

	p.setShield(p.shield - 1)
	mixer.playSound("hit")

	if p.shield <= 0 {
		p.beDestroyed()
	}
	deb.set(1, fmt.Sprintf("shield %d lives %d", p.shield, p.lives))
}

func (p *player) beDestroyed() {

	mixer.playSound("explosion")
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
	p.power = 0
	p.helperCount = 0
}

func (p *player) executeCollisionWith(other entity) {

	if other.getType() == entityTypeEnemyBullet {
		p.beHit()
	} else if other.getType() == entityTypeEnemyBig {
		p.beDestroyed()
	} else if other.getType() == entityTypeEnemySmall {
		p.beDestroyed()
	} else if other.getType() == entityTypeEnemyExtra {
		p.beDestroyed()
	} else if other.getType() == entityTypeBomb {
		p.beDestroyed()
	} else if other.getType() == entityTypePowerUpBullet {
		p.powerUp(other)
	} else if other.getType() == entityTypePowerUpHelper {
		p.powerUp(other)
	}
}

func (p *player) powerUp(powerUpEntity entity) {

	if powerUpEntity.getType() == entityTypePowerUpBullet {
		if p.power < playerMaxPower {
			p.power++
		}
	} else if powerUpEntity.getType() == entityTypePowerUpHelper {
		if p.helperCount < playerMaxHelper {
			p.helperCount++
		}
	}
	score.incrementPointsP1(5)
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
