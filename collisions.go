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

func checkCollisions(enemiesBig []*enemyBig, enemiesSmall []*enemySmall) {

	checkCollisionsEnemiesBig(enemiesBig)
	checkCollisionsEnemiesSmall(enemiesSmall)
	checkCollisionsPlayer(enemiesBig)
}

func checkCollisionsEnemiesBig(enemiesBig []*enemyBig) {

	for _, en := range enemiesBig {
		if en.active {
			for _, bul := range bulletPool {
				if bul.active && bul.fromPlayer {
					cEn := circle{
						x:      en.x,
						y:      en.y,
						radius: enemyBigCollisionRadius,
					}
					cBul := circle{
						x:      bul.x,
						y:      bul.y,
						radius: bulletCollisionRadius,
					}
					if collides(cEn, cBul) {
						bul.active = false
						en.beHit()
					}
				}
			}
		}
	}
}

func checkCollisionsEnemiesSmall(enemiesSmall []*enemySmall) {

	for _, en := range enemiesSmall {
		if en.active {
			for _, bul := range bulletPool {
				if bul.active && bul.fromPlayer {
					cEn := circle{
						x:      en.x,
						y:      en.y,
						radius: enemySmallCollisionRadius,
					}
					cBul := circle{
						x:      bul.x,
						y:      bul.y,
						radius: bulletCollisionRadius,
					}
					if collides(cEn, cBul) {
						bul.active = false
						en.beHit()
					}
				}
			}
		}
	}
}

func checkCollisionsPlayer(enemiesBig []*enemyBig) {

	if !plr.active {
		return
	}

	cPlr := circle{
		x:      plr.x,
		y:      plr.y,
		radius: playerCollisionRadius,
	}

	for _, en := range enemiesBig {
		if en.active {
			cEn := circle{
				x:      en.x,
				y:      en.y,
				radius: enemyBigCollisionRadius,
			}
			if collides(cEn, cPlr) {
				en.beDestroyed()
				plr.beDestroyed()
			}
		}
	}

	for _, en := range enemiesSmall {
		if en.active {
			cEn := circle{
				x:      en.x,
				y:      en.y,
				radius: enemySmallCollisionRadius,
			}
			if collides(cEn, cPlr) {
				en.beDestroyed()
				plr.beDestroyed()
			}
		}
	}

	for _, bul := range bulletPool {
		if bul.active && !bul.fromPlayer {
			cBul := circle{
				x:      bul.x,
				y:      bul.y,
				radius: bulletCollisionRadius,
			}
			if collides(cPlr, cBul) {
				bul.active = false
				plr.beHit()
			}
		}
	}
}
