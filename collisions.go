package main

import "math"

type circle struct {
	x, y   float64
	radius float64
}

func collides(c1, c2 circle) bool {
	dist := math.Sqrt(math.Pow(c2.x-c1.x, 2) + math.Pow(c2.y-c1.y, 2))

	return dist <= c1.radius+c2.radius
}

func checkCollisions(enemiesBig []entity, enemiesSmall []entity) {

	checkCollisionsEnemies(enemiesBig)
	checkCollisionsEnemies(enemiesSmall)
	checkCollisionsPlayer()
}

func checkCollisionsEnemies(enemiesBig []entity) {

	for _, en := range enemiesBig {
		if en.isActive() {
			for _, bul := range bulletPool {
				if bul.isActive() && collides(en.getCollisionCircle(), bul.getCollisionCircle()) {
					en.executeCollisionWith(bul)
					bul.executeCollisionWith(en)
				}
			}
		}
	}
}

func checkCollisionsPlayer() {

	if !plr.isActive() {
		return
	}

	cPlr := plr.getCollisionCircle()

	for _, en := range enemiesBig {
		if en.isActive() {
			if collides(en.getCollisionCircle(), cPlr) {
				en.executeCollisionWith(plr)
				plr.executeCollisionWith(en)
			}
		}
	}

	for _, en := range enemiesSmall {
		if en.isActive() {
			if collides(en.getCollisionCircle(), cPlr) {
				en.executeCollisionWith(plr)
				plr.executeCollisionWith(en)
			}
		}
	}

	for _, bul := range bulletPool {
		if bul.isActive() && bul.getType() == entityTypeEnemyBullet {
			if collides(cPlr, bul.getCollisionCircle()) {
				bul.executeCollisionWith(plr)
				plr.executeCollisionWith(bul)
			}
		}
	}

	for _, pu := range powerUps {
		if pu.isActive() {
			if collides(pu.getCollisionCircle(), cPlr) {
				pu.executeCollisionWith(plr)
				plr.executeCollisionWith(pu)
			}
		}
	}
}
