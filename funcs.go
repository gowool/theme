package theme

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"html"
	"html/template"
	"reflect"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/davecgh/go-spew/spew"
	"github.com/segmentio/go-camelcase"
	"github.com/segmentio/go-snakecase"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"

	"github.com/gowool/theme/internal"
)

var Funcs = template.FuncMap{
	"empty": func(given any) bool {
		g := reflect.ValueOf(given)
		return !g.IsValid() || g.IsNil() || g.IsZero()
	},
	"raw": func(s string) template.HTML {
		return template.HTML(s)
	},
	"escape": html.EscapeString,
	"deref": func(s any) any {
		v := reflect.ValueOf(s)
		if v.Kind() == reflect.Pointer {
			return v.Elem().Interface()
		}
		return s
	},
	"dump":           spew.Sdump,
	"str_camelcase":  camelcase.Camelcase,
	"str_snakecase":  snakecase.Snakecase,
	"str_trim":       strings.TrimSpace,
	"str_trim_left":  strings.TrimPrefix,
	"str_trim_right": strings.TrimSuffix,
	"str_upper":      strings.ToUpper,
	"str_lower":      strings.ToLower,
	"str_title":      strings.ToTitle,
	"str_contains":   strings.Contains,
	"str_has_prefix": strings.HasPrefix,
	"str_has_suffix": strings.HasSuffix,
	"str_replace":    strings.ReplaceAll,
	"str_equal":      strings.EqualFold,
	"str_index":      strings.Index,
	"str_join":       strings.Join,
	"str_split":      strings.Split,
	"str_split_n":    strings.SplitN,
	"str_fields":     strings.Fields,
	"str_repeat":     strings.Repeat,
	"str_len":        func(s string) int { return utf8.RuneCountInString(s) },
	"json": func(v any) string {
		return encode(v, json.Marshal)
	},
	"xml": func(v any) string {
		return encode(v, xml.Marshal)
	},
	"yaml": func(v any) string {
		return encode(v, yaml.Marshal)
	},
	"json_pretty": func(v any) string {
		return pretty(v, json.MarshalIndent)
	},
	"xml_pretty": func(v any) string {
		return pretty(v, xml.MarshalIndent)
	},
	"yaml_pretty": func(v any) string {
		var buf bytes.Buffer
		enc := yaml.NewEncoder(&buf)
		enc.SetIndent(2)
		if err := enc.Encode(v); err != nil {
			return ""
		}
		return template.JSEscapeString(internal.String(buf.Bytes()))
	},
	"len":   func(v any) int { return reflect.ValueOf(v).Len() },
	"list":  func(v ...any) []any { return v },
	"at":    func(v []any, i int) any { return v[i] },
	"first": func(v []any) any { return v[0] },
	"last":  func(v []any) any { return v[len(v)-1] },
	"slice": func(v []any, indices ...int) []any {
		switch len(indices) {
		case 0:
			return v[:]
		case 1:
			return v[indices[0]:]
		default:
			return v[indices[0]:indices[1]]
		}
	},
	"append":   func(v []any, e ...any) []any { return append(v, e...) },
	"prepend":  func(v []any, e ...any) []any { return append(e, v...) },
	"reverse":  func(v []any) []any { slices.Reverse(v); return v },
	"repeat":   func(v []any, count int) []any { return slices.Repeat(v, count) },
	"contains": func(v []any, i any) bool { return slices.Contains(v, i) },
	"index":    func(v []any, i any) int { return slices.Index(v, i) },
	"concat":   func(sl ...[]any) []any { return slices.Concat(sl...) },
	"dict": func(v ...any) map[any]any {
		if len(v)%2 != 0 {
			v = append(v, "")
		}
		dict := make(map[any]any, len(v)/2)
		for i := 0; i < len(v); i += 2 {
			dict[v[i]] = v[i+1]
		}
		return dict
	},
	"keys":   func(m map[any]any) []any { return maps.Keys(m) },
	"values": func(m map[any]any) []any { return maps.Values(m) },
	"has":    func(m map[any]any, k any) bool { _, ok := m[k]; return ok },
	"get":    func(m map[any]any, k any) any { return m[k] },
	"set":    func(m map[any]any, k, v any) map[any]any { m[k] = v; return m },
	"unset":  func(m map[any]any, k any) map[any]any { delete(m, k); return m },
	"now":    time.Now,
	"date":   formatDate,
	"date_local": func(fmt string, date any) string {
		return formatDate(fmt, date, "Local")
	},
	"date_utc": func(fmt string, date any) string {
		return formatDate(fmt, date, "UTC")
	},
}

func formatDate(fmt string, date any, location string) string {
	var t time.Time
	switch date := date.(type) {
	case time.Time:
		t = date
	case *time.Time:
		t = *date
	case int64:
		t = time.Unix(date, 0)
	case int:
		t = time.Unix(int64(date), 0)
	case int32:
		t = time.Unix(int64(date), 0)
	default:
		t = time.Now()
	}

	loc, err := time.LoadLocation(location)
	if err != nil {
		loc, _ = time.LoadLocation("UTC")
	}

	return t.In(loc).Format(fmt)
}

func encode(v any, fn func(v any) ([]byte, error)) string {
	raw, err := fn(v)
	if err != nil {
		return ""
	}
	return template.JSEscapeString(internal.String(raw))
}

func pretty(v any, fn func(v any, prefix, indent string) ([]byte, error)) string {
	raw, err := fn(v, "", "  ")
	if err != nil {
		return ""
	}
	return template.JSEscapeString(internal.String(raw))
}
