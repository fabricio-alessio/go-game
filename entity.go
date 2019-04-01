package main

const (
	entityTypeEnemyBig     = 0
	entityTypeEnemySmall   = 1
	entityTypePlayer       = 2
	entityTypePlayerBullet = 3
	entityTypeEnemyBullet  = 4
)

type entity interface {
	executeCollisionWith(other entity)
	draw()
	update()
	getType() int8
	deactivate()
	isActive() bool
	start(x, y, angle, speed float64, entityType int8)
	getCollisionCircle() circle
}
