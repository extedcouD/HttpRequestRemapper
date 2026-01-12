package httprequestremap

import (
	"strings"
	"time"

	"github.com/AsaiYusuke/jsonpath"
)

type Builtins struct {
	UUID string
	Now  func() string
}

func ApplyTemplate(root any, template any, builtins Builtins) any {
	switch t := template.(type) {
	case map[string]any:
		out := map[string]any{}
		for k, v := range t {
			out[k] = ApplyTemplate(root, v, builtins)
		}
		return out
	case []any:
		out := make([]any, 0, len(t))
		for _, v := range t {
			out = append(out, ApplyTemplate(root, v, builtins))
		}
		return out
	case string:
		expr := strings.TrimSpace(t)
		if expr == "" {
			return ""
		}

		if expr == "uuid()" {
			return builtins.UUID
		}
		if expr == "now()" {
			if builtins.Now != nil {
				return builtins.Now()
			}
			return time.Now().UTC().Format(time.RFC3339Nano)
		}

		if strings.HasPrefix(expr, "$") {
			return EvalJSONPath(root, expr)
		}
		return t
	default:
		return template
	}
}

func EvalJSONPath(root any, expr string) any {
	results, err := jsonpath.Retrieve(expr, root)
	if err != nil {
		return nil
	}
	if len(results) == 0 {
		return nil
	}
	if len(results) == 1 {
		return results[0]
	}
	return results
}
