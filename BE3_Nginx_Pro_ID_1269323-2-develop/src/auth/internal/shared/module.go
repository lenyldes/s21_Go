package shared

import (
	"go.uber.org/fx"
)

const module_name = "shared"

func Module() fx.Option {
	return fx.Module(
		module_name,
		fx.Provide(NewConfig),
		fx.Provide(NewLogger),
	)
}
