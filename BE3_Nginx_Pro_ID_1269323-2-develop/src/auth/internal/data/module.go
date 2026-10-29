package data

import (
	"auth/internal/auth"

	"go.uber.org/fx"
)

const MODULE_NAME = "data"

func Module() fx.Option {
	return fx.Module(
		"data",
		fx.Provide(
			fx.Annotate(
				NewUserRepo,
				fx.As(new(auth.AuthRepository)),
			),
			fx.Annotate(
				NewUserRepo,
				fx.As(new(auth.UserProvider)),
			),
		),
		fx.Provide(InitPostgresPool),
	)
}
