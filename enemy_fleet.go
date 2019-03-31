package main

import (
	"math/rand"
	"time"
)

const (
	fleetReleaseBigCooldown       = time.Millisecond * 4677
	fleetTimeToExecuteRestarLevel = time.Millisecond * 3000
	fleetReleaseSmallCooldown     = time.Millisecond * 250
	fleetReleaseSmallFaseCooldown = time.Millisecond * 13452
	fleetReleaseSmallMax          = 10
)

type enemyFleet struct {
	lastTimeReleaseBig       time.Time
	timeAskRestartLevel      time.Time
	timeAskGoToMenu          time.Time
	lastTimeReleaseSmall     time.Time
	lastTimeReleaseSmallFase time.Time
	actualLevel              int8
	mustRestartLevel         bool
	mustGoToMenu             bool
	smallFaseX               float64
	smallReleased            int8
	smallFase                int8
}

func newEnemyFleet() (ef enemyFleet) {

	return ef
}

func (ef *enemyFleet) update() {

	if time.Since(ef.lastTimeReleaseBig) >= fleetReleaseBigCooldown {
		ef.releaseBig()

		ef.lastTimeReleaseBig = time.Now()
	}

	if time.Since(ef.lastTimeReleaseSmallFase) >= fleetReleaseSmallFaseCooldown {
		ef.releaseSmallFase()

		ef.lastTimeReleaseSmallFase = time.Now()
	}

	if time.Since(ef.lastTimeReleaseSmall) >= fleetReleaseSmallCooldown {
		ef.releaseSmall()

		ef.lastTimeReleaseSmall = time.Now()
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

	en := enemyBigFromPool()
	if en != nil {
		en.x = float64(rand.Intn(screenWidth))
		en.y = -30
		en.hitCount = 0
		en.active = true
	}
}

func (ef *enemyFleet) releaseSmallFase() {

	ef.smallFaseX = float64(rand.Intn(screenWidth-400) + 400)
	ef.smallReleased = 0
	ef.smallFase++
}

func (ef *enemyFleet) releaseSmall() {

	if ef.smallReleased >= fleetReleaseSmallMax {
		return
	}

	en := enemySmallFromPool()
	if en != nil {
		en.start(ef.smallReleased, ef.smallFase, ef.smallFaseX, -30)
		ef.smallReleased++
	}
}

func (ef *enemyFleet) startLevel(level int8) {

	ef.actualLevel = level
	ef.smallFase = 0
}

func (ef *enemyFleet) askRestartLevel() {

	ef.mustRestartLevel = true
	ef.timeAskRestartLevel = time.Now()
}

func (ef *enemyFleet) restartLevel() {
	deactivateAllBullets()
	deactivateAllEnemiesBig()
	ef.smallFase = 0
	ef.mustRestartLevel = false
	plr.start()
}

func (ef *enemyFleet) askGoToMenu() {

	ef.mustGoToMenu = true
	ef.timeAskGoToMenu = time.Now()
}

func (ef *enemyFleet) goToMenu() {
	deactivateAllBullets()
	deactivateAllEnemiesBig()
	ef.smallFase = 0
	ef.mustGoToMenu = false
	gameStarted = false
}
