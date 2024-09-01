package theme

import (
	"context"
	"time"
)

type Template interface {
	Code() []byte
	Changed() time.Time
}

type Repository interface {
	FindByName(ctx context.Context, name string) (Template, error)
}
