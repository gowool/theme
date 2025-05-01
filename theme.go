package theme

import (
	"context"
	"html/template"
	"io"
	"strings"

	et "github.com/gowool/extends-template"
)

type (
	Handler          = et.Handler
	Source           = et.Source
	Loader           = et.Loader
	ChainLoader      = et.ChainLoader
	MemoryLoader     = et.MemoryLoader
	FileSystemLoader = et.FileSystemLoader

	FuncMap func(t Theme) template.FuncMap

	environment struct {
		*et.Environment
	}
)

var (
	NewChainLoader      = et.NewChainLoader
	NewMemoryLoader     = et.NewMemoryLoader
	NewFileSystemLoader = et.NewFileSystemLoader
	NewFSLoaderWithNS   = et.NewFSLoaderWithNS
)

type Theme interface {
	Debug(debug bool) Theme
	Funcs(funcMap FuncMap) Theme
	Global(global ...string) Theme
	Write(ctx context.Context, w io.Writer, name string, data any) error
	HTML(ctx context.Context, name string, data any) (string, error)
}

func New(loader Loader, handlers ...Handler) Theme {
	env := &environment{Environment: et.NewEnvironment(loader, handlers...)}
	env.Environment.Funcs(Funcs)
	return env
}

func (t *environment) Debug(debug bool) Theme {
	t.Environment.Debug(debug)
	return t
}

func (t *environment) Funcs(funcMap FuncMap) Theme {
	t.Environment.Funcs(funcMap(t))
	return t
}

func (t *environment) Global(global ...string) Theme {
	t.Environment.Global(global...)
	return t
}

func (t *environment) Write(ctx context.Context, w io.Writer, name string, data any) error {
	wrap, err := t.Load(ctx, name)
	if err != nil {
		return err
	}

	return wrap.HTML.ExecuteTemplate(w, name, data)
}

func (t *environment) HTML(ctx context.Context, name string, data any) (string, error) {
	var b strings.Builder
	if err := t.Write(ctx, &b, name, data); err != nil {
		return "", err
	}
	return b.String(), nil
}
