package theme

import (
	"time"

	"github.com/google/uuid"

	"github.com/gowool/theme/internal/util"
)

type Template struct {
	ID        uuid.UUID  `cfg:"id,omitempty" json:"id,omitempty" yaml:"id,omitempty" bson:"id,omitempty"`
	Code      string     `cfg:"code,omitempty" json:"code,omitempty" yaml:"code,omitempty" bson:"code,omitempty"`
	Content   string     `cfg:"content,omitempty" json:"content,omitempty" yaml:"content,omitempty" bson:"content,omitempty"`
	Created   time.Time  `cfg:"created,omitempty" json:"created,omitempty" yaml:"created,omitempty" bson:"created,omitempty"`
	Updated   time.Time  `cfg:"updated,omitempty" json:"updated,omitempty" yaml:"updated,omitempty" bson:"updated,omitempty"`
	Published *time.Time `cfg:"published,omitempty" json:"published,omitempty" yaml:"published,omitempty" bson:"published,omitempty"`
	Expired   *time.Time `cfg:"expired,omitempty" json:"expired,omitempty" yaml:"expired,omitempty" bson:"expired,omitempty"`
}

func (t *Template) String() string {
	return t.Code
}

func (t *Template) Enabled(now time.Time) bool {
	return t.Published != nil &&
		!t.Published.IsZero() &&
		(t.Published.Before(now) || t.Published.Equal(now)) &&
		(t.Expired == nil || t.Expired.IsZero() || t.Expired.After(now))
}

func (t *Template) ContentBytes() []byte {
	return util.StringToBytes(t.Content)
}
