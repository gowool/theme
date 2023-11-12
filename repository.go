package theme

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/gowool/cr"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, tags ...string) error
	Get(ctx context.Context, key string, value interface{}) error
	DelByKey(ctx context.Context, key string) error
	DelByTag(ctx context.Context, tag string) error
}

type Repository interface {
	Find(ctx context.Context, criteria *cr.Criteria) ([]*Template, int, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Template, error)
	FindByCode(ctx context.Context, code string, date time.Time) (*Template, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	Save(ctx context.Context, template *Template) error
}

type CacheRepository struct {
	Repository
	Cache Cache
}

func (r CacheRepository) FindByID(ctx context.Context, id uuid.UUID) (template *Template, err error) {
	key := "theme:template:id:" + id.String()

	if err = r.Cache.Get(ctx, key, &template); err == nil {
		return
	}

	if template, err = r.Repository.FindByID(ctx, id); err != nil {
		return
	}

	_ = r.Cache.Set(ctx, key, template, "theme:template:tag:"+id.String())

	return
}

func (r CacheRepository) FindByCode(ctx context.Context, code string, date time.Time) (template *Template, err error) {
	key := "theme:template:code:" + code

	if err = r.Cache.Get(ctx, key, &template); err == nil {
		if template.Enabled(date) {
			return
		}

		_ = r.Cache.DelByKey(ctx, key)
	}

	if template, err = r.Repository.FindByCode(ctx, code, date); err != nil {
		return
	}

	_ = r.Cache.Set(ctx, key, template, "theme:template:tag:"+template.ID.String())

	return
}

func (r CacheRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	defer r.del(ctx, id)

	return r.Repository.DeleteByID(ctx, id)
}

func (r CacheRepository) Save(ctx context.Context, template *Template) error {
	defer func() { r.del(ctx, template.ID) }()

	return r.Repository.Save(ctx, template)
}

func (r CacheRepository) del(ctx context.Context, id uuid.UUID) {
	_ = r.Cache.DelByTag(ctx, "theme:template:tag:"+id.String())
}
