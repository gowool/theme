package fx

import (
	"go.uber.org/fx"

	"github.com/gowool/theme"
)

var (
	OptionTheme  = fx.Provide(NewTheme)
	OptionLoader = fx.Provide(
		fx.Annotate(
			theme.NewRepositoryLoader,
			fx.As(new(theme.Loader)),
		),
	)
)
