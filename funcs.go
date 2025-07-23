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
	"ternary": func(condition bool, trueValue, falseValue any) any {
		if condition {
			return trueValue
		}
		return falseValue
	},
	"js":        func(str string) template.JS { return template.JS(str) },
	"css":       func(str string) template.CSS { return template.CSS(str) },
	"raw":       func(s string) template.HTML { return template.HTML(s) },
	"html":      func(str string) template.HTML { return template.HTML(str) },
	"html_attr": func(str string) template.HTMLAttr { return template.HTMLAttr(str) },
	"empty": func(given any) bool {
		g := reflect.ValueOf(given)
		return !g.IsValid() || g.IsNil() || g.IsZero()
	},
	"escape": html.EscapeString,
	"deref": func(s any) any {
		v := reflect.ValueOf(s)
		if v.Kind() == reflect.Pointer {
			return v.Elem().Interface()
		}
		return s
	},
	"dump": spew.Sdump,
	"str_build": func(str ...string) string {
		var b strings.Builder
		for _, s := range str {
			b.WriteString(s)
		}
		return b.String()
	},
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
		return Encode(v, json.Marshal)
	},
	"xml": func(v any) string {
		return Encode(v, xml.Marshal)
	},
	"yaml": func(v any) string {
		return Encode(v, yaml.Marshal)
	},
	"json_pretty": func(v any) string {
		return Pretty(v, json.MarshalIndent)
	},
	"xml_pretty": func(v any) string {
		return Pretty(v, xml.MarshalIndent)
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
	"seq":   Seq,
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
	"index_of": func(v []any, i any) int { return slices.Index(v, i) },
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
	"date":   FormatDate,
	"date_local": func(fmt string, date any) string {
		return FormatDate(fmt, date, "Local")
	},
	"date_utc": func(fmt string, date any) string {
		return FormatDate(fmt, date, "UTC")
	},
}

func FormatDate(fmt string, date any, location string) string {
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

func Encode(v any, fn func(v any) ([]byte, error)) string {
	raw, err := fn(v)
	if err != nil {
		return ""
	}
	return template.JSEscapeString(internal.String(raw))
}

func Pretty(v any, fn func(v any, prefix, indent string) ([]byte, error)) string {
	raw, err := fn(v, "", "  ")
	if err != nil {
		return ""
	}
	return template.JSEscapeString(internal.String(raw))
}

// Seq creates a sequence of integers from args.
//
// Examples:
//
//	3 => 1, 2, 3
//	1 2 4 => 1, 3
//	-3 => -1, -2, -3
//	1 4 => 1, 2, 3, 4
//	1 -2 => 1, 0, -1, -2
func Seq(args ...int) []int {
	if len(args) < 1 || len(args) > 3 {
		// invalid number of arguments to Seq
		return nil
	}

	inc := 1
	var last int
	first := args[0]

	if len(args) == 1 {
		last = first
		if last == 0 {
			return nil
		} else if last > 0 {
			first = 1
		} else {
			first = -1
			inc = -1
		}
	} else if len(args) == 2 {
		last = args[1]
		if last < first {
			inc = -1
		}
	} else {
		inc = args[1]
		last = args[2]
		if inc == 0 {
			// 'increment' must not be 0
			return nil
		}
		if first < last && inc < 0 {
			// 'increment' must be > 0
			return nil
		}
		if first > last && inc > 0 {
			// 'increment' must be < 0
			return nil
		}
	}

	// sanity check
	if last < -100000 {
		// size of result exceeds limit
		return nil
	}
	size := ((last - first) / inc) + 1

	// sanity check
	if size <= 0 || size > 2000 {
		// size of result exceeds limit
		return nil
	}

	seq := make([]int, size)
	val := first
	for i := 0; ; i++ {
		seq[i] = val
		val += inc
		if (inc < 0 && val < last) || (inc > 0 && val > last) {
			break
		}
	}

	return seq
}
