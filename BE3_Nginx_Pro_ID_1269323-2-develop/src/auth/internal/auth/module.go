package auth

import (
	"go.uber.org/fx"
)

const module_name = "auth"

func Module() fx.Option {
	return fx.Module(
		module_name,
		fx.Provide(NewJwtCache),
		fx.Provide(NewKeys),
		fx.Provide(NewJwtProvider),
		fx.Provide(NewUserAuthenticator),
		fx.Provide(NewChangePasswordNotifier),
		fx.Provide(NewAuthService),
		fx.Provide(Options),
	)
}
