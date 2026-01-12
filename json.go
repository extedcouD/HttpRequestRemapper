package httprequestremap

import (
	"bytes"
	"encoding/json"
)

func TryParseJSON(b []byte) (any, bool) {
	if len(bytes.TrimSpace(b)) == 0 {
		return nil, false
	}
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, false
	}
	return v, true
}

func ParseJSONObjectOrEmpty(b []byte) map[string]any {
	v, ok := TryParseJSON(b)
	if !ok {
		return map[string]any{}
	}
	m, ok := v.(map[string]any)
	if !ok || m == nil {
		return map[string]any{}
	}
	return m
}
