package cache

import (
	"github.com/Echin-h/HangZhou-Monopoly/config"
	"github.com/Echin-h/HangZhou-Monopoly/internal/core/cache/driver"
	"github.com/Echin-h/HangZhou-Monopoly/internal/core/cache/types"
)

type Creator interface {
	Create(conf config.Cache) (types.Cache, error)
}

func init() {
	typeMap["redis"] = driver.RedisCreator{}
}

var typeMap = make(map[string]Creator)

func getCreatorByType(cacheType string) Creator {
	return typeMap[cacheType]
}
