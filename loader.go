package theme

import (
	"context"
	"fmt"

	"github.com/gowool/theme/internal"
)

var _ Loader = (*RepositoryLoader)(nil)

type RepositoryLoader struct {
	repository Repository
}

func NewRepositoryLoader(repository Repository) *RepositoryLoader {
	return &RepositoryLoader{repository: repository}
}

func (l *RepositoryLoader) Get(ctx context.Context, name string) (*Source, error) {
	if template, err := l.repository.FindByName(ctx, name); err == nil {
		return &Source{
			Name: name,
			Code: template.Code(),
		}, nil
	}

	// if the template is not available, it should return a dummy source.
	return &Source{
		Name: name,
		Code: internal.Bytes(fmt.Sprintf("{{raw \"<!-- ====== template `%s` is disabled ====== -->\"}}", name)),
	}, nil
}

func (l *RepositoryLoader) IsFresh(ctx context.Context, name string, t int64) (bool, error) {
	item, err := l.repository.FindByName(ctx, name)
	if err != nil {
		return false, err
	}
	return item.Changed().Unix() < t, nil
}

func (l *RepositoryLoader) Exists(ctx context.Context, name string) (bool, error) {
	_, err := l.repository.FindByName(ctx, name)
	return err == nil, err
}
