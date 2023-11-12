package theme

import (
	"context"
	"html/template"
	"io"
	"strings"

	"github.com/Masterminds/sprig/v3"
	et "github.com/gowool/extends-template"
)

var _ Theme = (*Environment)(nil)

type (
	FuncMap          = template.FuncMap
	Handler          = et.Handler
	Source           = et.Source
	Loader           = et.Loader
	ChainLoader      = et.ChainLoader
	MemoryLoader     = et.MemoryLoader
	FilesystemLoader = et.FilesystemLoader

	BeforeWrite func(ctx context.Context, name string, data interface{}) (newName string, newData interface{}, err error)

	Theme interface {
		Debug(debug bool) Theme
		Funcs(funcMap FuncMap) Theme
		Global(global ...string) Theme
		Write(ctx context.Context, w io.Writer, name string, data interface{}) error
		HTML(ctx context.Context, name string, data interface{}) (string, error)
	}

	Environment struct {
		*et.Environment
		beforeWrite BeforeWrite
	}
)

func New(loader Loader, beforeWrite BeforeWrite, handlers ...Handler) *Environment {
	env := &Environment{
		Environment: et.NewEnvironment(loader, handlers...),
		beforeWrite: beforeWrite,
	}

	funcMap := sprig.FuncMap()
	funcMap["raw"] = func(s string) template.HTML {
		return template.HTML(s)
	}

	env.Funcs(funcMap)

	return env
}

func (t *Environment) Debug(debug bool) Theme {
	t.Environment.Debug(debug)
	return t
}

func (t *Environment) Funcs(funcMap template.FuncMap) Theme {
	t.Environment.Funcs(funcMap)
	return t
}

func (t *Environment) Global(global ...string) Theme {
	t.Environment.Global(global...)
	return t
}

func (t *Environment) Write(ctx context.Context, w io.Writer, name string, data interface{}) error {
	if t.beforeWrite != nil {
		var err error
		if name, data, err = t.beforeWrite(ctx, name, data); err != nil {
			return err
		}
	}

	wrap, err := t.Environment.Load(ctx, name)
	if err != nil {
		return err
	}

	return wrap.HTML.ExecuteTemplate(w, name, data)
}

func (t *Environment) HTML(ctx context.Context, name string, data interface{}) (string, error) {
	var b strings.Builder
	if err := t.Write(ctx, &b, name, data); err != nil {
		return "", err
	}
	return b.String(), nil
}
