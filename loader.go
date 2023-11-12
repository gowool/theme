package theme

import (
	"context"
	"fmt"
	"time"

	"github.com/gowool/theme/internal/util"
)

var _ Loader = TemplateLoader{}

type TemplateLoader struct {
	Repository Repository
}

func (s TemplateLoader) Get(ctx context.Context, name string) (*Source, error) {
	if item, err := s.find(ctx, name); err == nil {
		return &Source{
			Name: name,
			Code: util.StringToBytes(item.Content),
		}, nil
	}

	// if the template is expired or not published, it should return a dummy source.
	return &Source{
		Name: name,
		Code: util.StringToBytes(fmt.Sprintf("{{raw \"<!-- ====== template `%s` is disabled ====== -->\"}}", name)),
	}, nil
}

func (s TemplateLoader) IsFresh(ctx context.Context, name string, t int64) (bool, error) {
	item, err := s.find(ctx, name)
	if err != nil {
		return false, err
	}

	return item.Updated.Unix() < t, nil
}

func (s TemplateLoader) Exists(ctx context.Context, name string) (bool, error) {
	_, err := s.find(ctx, name)
	return err == nil, err
}

func (s TemplateLoader) find(ctx context.Context, name string) (*Template, error) {
	return s.Repository.FindByCode(ctx, name, time.Now())
}
