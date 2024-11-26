package fx

import "go.uber.org/fx"

const ModuleName = "theme"

var Module = fx.Module(
	ModuleName,
	OptionTheme,
)
