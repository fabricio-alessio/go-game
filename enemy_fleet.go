package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	fleetReleaseBigCooldown       = time.Millisecond * 4677
	fleetTimeToExecuteRestarLevel = time.Millisecond * 3000
	fleetReleaseSmallCooldown     = time.Millisecond * 250
	fleetReleaseSmallMax          = 10
)

type attack struct {
	distance  float64
	quantity  int
	enemyType int8
}

type enemyFleet struct {
	lastTimeReleaseBig   time.Time
	timeAskRestartLevel  time.Time
	timeAskGoToMenu      time.Time
	lastTimeReleaseSmall time.Time
	actualLevel          int8
	mustRestartLevel     bool
	mustGoToMenu         bool
	level                int8
	position             float64
	fleetAttacks         []attack
	attackNumber         int
}

func newEnemyFleet() (ef enemyFleet) {

	ef.fleetAttacks = append(ef.fleetAttacks, attack{distance: 500, quantity: 1, enemyType: entityTypeEnemyBig})
	ef.fleetAttacks = append(ef.fleetAttacks, attack{distance: 500, quantity: 2, enemyType: entityTypeEnemySmall})
	ef.fleetAttacks = append(ef.fleetAttacks, attack{distance: 500, quantity: 2, enemyType: entityTypeEnemyExtraComming})
	ef.fleetAttacks = append(ef.fleetAttacks, attack{distance: 500, quantity: 2, enemyType: entityTypeEnemyBig})
	return ef
}

func (ef *enemyFleet) update() {

	ef.position += delta
	someAttack := ef.fleetAttacks[ef.attackNumber]
	deb.set(5, fmt.Sprintf("position %f", ef.position))

	if ef.position > someAttack.distance {
		ef.releaseAttack(someAttack)
		ef.position = 0
		ef.prepareNextAttack()
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

func (ef *enemyFleet) prepareNextAttack() {

	if ef.attackNumber+1 >= len(ef.fleetAttacks) {
		ef.attackNumber = 0
	} else {
		ef.attackNumber++
	}
}

func (ef *enemyFleet) releaseAttack(someAttack attack) {

	fmt.Printf("Attack %s released %d times, sequency %d\n", getNameOfType(someAttack.enemyType), someAttack.quantity, ef.attackNumber)
	for i := 0; i < someAttack.quantity; i++ {
		if someAttack.enemyType == entityTypeEnemyBig {
			ef.releaseBig()
		} else if someAttack.enemyType == entityTypeEnemySmall {
			ef.releaseSmall()
		} else if someAttack.enemyType == entityTypeEnemyExtraComming {
			ef.releaseExtraComming()
		}
	}
}

func (ef *enemyFleet) releaseBig() {

	en := enemyBigFromPool()
	if en != nil {
		// TODO send data
		en.start(0, 0, 0, 0, 0)
	}
}

func (ef *enemyFleet) releaseExtraComming() {

	en := enemyExtraCommingFromPool()
	if en != nil {
		en.start(0, 0, 0, 0, 0)
	}
}

func (ef *enemyFleet) releaseSmall() {

	x := float64(rand.Intn(int(screenWidth)-400) + 400)
	group := enemyGroup{
		aliveCount:     6,
		variationIndex: rand.Intn(enemySmallValCount)}

	for i := 0; i < 6; i++ {
		en := enemySmallFromPool()
		if en != nil {
			enSmall, ok := en.(*enemySmall)
			if ok {
				enSmall.group = &group
			}
			en.start(x, float64(-30-(i*40)), 0, 0, 0)
		}
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
	deactivateAllBulletBlasts()
	deactivateAllEnemiesBig()
	deactivateAllEnemiesSmall()
	deactivateAllEnemiesExtra()
	deactivateAllEnemiesExtraComming()
	ef.mustRestartLevel = false
	plr.start(0, 0, 0, 0, 0)
}

func (ef *enemyFleet) askGoToMenu() {

	ef.mustGoToMenu = true
	ef.timeAskGoToMenu = time.Now()
}

func (ef *enemyFleet) goToMenu() {
	deactivateAllBullets()
	deactivateAllEnemiesBig()
	ef.mustGoToMenu = false
	gameStarted = false
}
