package battle

import "github.com/Noblefel/term-rpg/internal/entity"

const (
	WIN int = iota
	LOSE
	FLED
	NEXT
)

type Battle struct {
	Turn        int
	IsEnemyTurn bool
	Log         string
	Status      int
	Enemy       entity.Enemy
	EnemyAttr   *entity.EnemyBase
}

func New(e entity.Enemy) *Battle {
	return &Battle{
		Log:       "-- No Recent Log --",
		Enemy:     e,
		EnemyAttr: e.Attr(),
	}
}
