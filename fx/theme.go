package fx

import (
	"go.uber.org/fx"

	"github.com/gowool/theme"
)

type ThemeParams struct {
	fx.In
	Config   Config
	Loader   theme.Loader
	Funcs    []theme.FuncMap `group:"theme-func-map"`
	Handlers []theme.Handler `group:"theme-handler"`
}

func NewTheme(params ThemeParams) theme.Theme {
	t := theme.New(params.Loader, params.Handlers...).Debug(params.Config.Debug)

	for _, f := range params.Funcs {
		t = t.Funcs(f)
	}

	return t
}