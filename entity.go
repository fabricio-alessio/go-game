package main

const (
	entityTypeEnemyBig          = 0
	entityTypeEnemySmall        = 1
	entityTypePlayer            = 2
	entityTypePlayerBullet      = 3
	entityTypeEnemyBullet       = 4
	entityTypePowerUp           = 5
	entityTypeEnemyExtra        = 6
	entityTypeEnemyExtraComming = 7
	entityTypeBulletBlast       = 8
)

func getNameOfType(entityType int8) string {

	switch entityType {
	case entityTypeEnemyBig:
		return "EnemyBig"
	case entityTypeEnemySmall:
		return "EnemySmall"
	case entityTypePlayer:
		return "Player"
	case entityTypePlayerBullet:
		return "PlayerBullet"
	case entityTypeEnemyBullet:
		return "EnemyBullet"
	case entityTypePowerUp:
		return "PowerUp"
	case entityTypeEnemyExtra:
		return "EnemyExtra"
	case entityTypeEnemyExtraComming:
		return "EnemyExtraComming"
	case entityTypeBulletBlast:
		return "BulletBlast"
	default:
		return "Undefined"
	}
}

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
