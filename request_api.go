package httprequestremap

import (
	"net/http"
)

type RequestOptions struct {
	// MaxBodyBytes controls how many bytes of the request body are captured.
	// If <= 0, a default of 1MiB is used.
	MaxBodyBytes int64
}

func (o RequestOptions) maxBodyBytesOrDefault() int64 {
	if o.MaxBodyBytes <= 0 {
		return 1024 * 1024
	}
	return o.MaxBodyBytes
}

// RootFromRequest captures the request body (restoring r.Body) and builds the
// canonical JSONPath root for a single *http.Request.
//
// The returned root contains:
// - method, path, host
// - headers, headers_all
// - query, query_all
// - cookies
// - body (parsed JSON if possible, else {})
func RootFromRequest(r *http.Request, opts *RequestOptions) (root map[string]any, present bool, body []byte, truncated bool) {
	var o RequestOptions
	if opts != nil {
		o = *opts
	}

	present, body, truncated = CaptureRequestBody(r, o.maxBodyBytesOrDefault())
	root = RequestRoot(r, body)
	return root, present, body, truncated
}

// EvalJSONPathFromRequest evaluates a JSONPath expression against the canonical
// root built from the given request.
func EvalJSONPathFromRequest(r *http.Request, expr string, opts *RequestOptions) any {
	root, _, _, _ := RootFromRequest(r, opts)
	return EvalJSONPath(root, expr)
}

// ApplyTemplateFromRequest applies a remap template against the canonical root
// built from the given request.
//
// Template semantics are the same as ApplyTemplate:
// - strings starting with '$' are JSONPath expressions
// - other strings are literals
// - builtins: uuid(), now()
func ApplyTemplateFromRequest(r *http.Request, template any, builtins Builtins, opts *RequestOptions) any {
	root, _, _, _ := RootFromRequest(r, opts)
	return ApplyTemplate(root, template, builtins)
}
