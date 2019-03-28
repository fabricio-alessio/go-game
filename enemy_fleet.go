package main

import (
	"math/rand"
	"time"
)

const (
	fleetReleaseBigCooldown       = time.Millisecond * 4677
	fleetTimeToExecuteRestarLevel = time.Millisecond * 3000
)

type enemyFleet struct {
	lastTimeReleaseBig  time.Time
	timeAskRestartLevel time.Time
	timeAskGoToMenu     time.Time
	actualLevel         int8
	mustRestartLevel    bool
	mustGoToMenu        bool
}

func newEnemyFleet() (ef enemyFleet) {

	return ef
}

func (ef *enemyFleet) update() {

	if time.Since(ef.lastTimeReleaseBig) >= fleetReleaseBigCooldown {
		ef.releaseBig()

		ef.lastTimeReleaseBig = time.Now()
	}

	if ef.mustRestartLevel {
		if time.Since(ef.timeAskRestartLevel) >= fleetTimeToExecuteRestarLevel {
			ef.restartLevel()
		}
	}

	if ef.mustGoToMenu {
		if time.Since(ef.timeAskGoToMenu) >= fleetTimeToExecuteRestarLevel {
			ef.goToMenu()
		}
	}
}

func (ef *enemyFleet) releaseBig() {

	en := enemyFromPool()
	if en != nil {
		en.x = float64(rand.Intn(screenWidth))
		en.y = -30
		en.hitCount = 0
		en.active = true
	}
}

func (ef *enemyFleet) startLevel(level int8) {

	ef.actualLevel = level
}

func (ef *enemyFleet) askRestartLevel() {

	ef.mustRestartLevel = true
	ef.timeAskRestartLevel = time.Now()
}

func (ef *enemyFleet) restartLevel() {
	deactivateAllBullets()
	deactivateAllEnemies()
	ef.mustRestartLevel = false
	plr.start()
}

func (ef *enemyFleet) askGoToMenu() {

	ef.mustGoToMenu = true
	ef.timeAskGoToMenu = time.Now()
}

func (ef *enemyFleet) goToMenu() {
	deactivateAllBullets()
	deactivateAllEnemies()
	ef.mustGoToMenu = false
	gameStarted = false
}

func enemyFromPool() *enemy {
	for _, en := range enemies {
		if !en.active {
			return en
		}
	}

	return nil
}
