package httprequestremap

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// CaptureRequestBody reads up to maxBytes (+1 to detect truncation), then restores r.Body.
//
// Return values:
// - present: whether r.Body was non-nil
// - body: captured bytes (possibly truncated)
// - truncated: whether the body exceeded maxBytes
func CaptureRequestBody(r *http.Request, maxBytes int64) (present bool, body []byte, truncated bool) {
	if r == nil || r.Body == nil {
		return false, nil, false
	}
	if maxBytes == 0 {
		return true, nil, true
	}

	limited := io.LimitReader(r.Body, maxBytes+1)
	b, err := io.ReadAll(limited)
	if err != nil {
		// Preserve prior behavior: treat read failure as "captured but empty".
		r.Body = io.NopCloser(bytes.NewReader(nil))
		return true, nil, false
	}

	truncated = int64(len(b)) > maxBytes
	if truncated {
		b = b[:maxBytes]
	}

	r.Body = io.NopCloser(bytes.NewReader(b))
	return true, b, truncated
}

func HeaderMaps(h http.Header) (map[string]any, map[string]any) {
	first := map[string]any{}
	all := map[string]any{}
	for k, vs := range h {
		lk := strings.ToLower(k)
		if len(vs) == 0 {
			continue
		}
		all[lk] = append([]string(nil), vs...)
		first[lk] = vs[0]
	}
	return first, all
}

func QueryMaps(u *url.URL) (map[string]any, map[string]any) {
	first := map[string]any{}
	all := map[string]any{}
	if u == nil {
		return first, all
	}
	q := u.Query()
	for k, vs := range q {
		if len(vs) == 0 {
			continue
		}
		all[k] = append([]string(nil), vs...)
		first[k] = vs[0]
	}
	return first, all
}

func CookieMap(r *http.Request) map[string]any {
	res := map[string]any{}
	if r == nil {
		return res
	}
	for _, c := range r.Cookies() {
		if c == nil {
			continue
		}
		res[c.Name] = c.Value
	}
	return res
}

// RequestRoot builds a canonical JSON object derived from an *http.Request.
// This is useful as a JSONPath root for remapping.
func RequestRoot(r *http.Request, bodyBytes []byte) map[string]any {
	headersFirst, headersAll := HeaderMaps(nil)
	queryFirst, queryAll := map[string]any{}, map[string]any{}
	cookies := CookieMap(r)
	method := ""
	path := ""
	host := ""
	if r != nil {
		headersFirst, headersAll = HeaderMaps(r.Header)
		queryFirst, queryAll = QueryMaps(r.URL)
		method = r.Method
		host = r.Host
		if r.URL != nil {
			path = r.URL.Path
		}
	}

	var parsed any
	if v, ok := TryParseJSON(bodyBytes); ok {
		parsed = v
	} else {
		parsed = map[string]any{}
	}

	return map[string]any{
		"method":      method,
		"path":        path,
		"host":        host,
		"headers":     headersFirst,
		"headers_all": headersAll,
		"query":       queryFirst,
		"query_all":   queryAll,
		"cookies":     cookies,
		"body":        parsed,
	}
}
