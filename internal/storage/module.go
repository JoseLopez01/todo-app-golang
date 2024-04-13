package storage

import "go.uber.org/fx"

var Module = fx.Module(
	"storage-module",
	fx.Provide(NewRedisClient),
)
