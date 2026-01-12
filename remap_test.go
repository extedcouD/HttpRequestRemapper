package httprequestremap

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApplyTemplate_JSONPathAndBuiltins(t *testing.T) {
	root := map[string]any{
		"headers": map[string]any{"x": "1"},
		"cookies": map[string]any{"sid": "abc"},
	}

	got := ApplyTemplate(root, map[string]any{
		"x":    "$.headers.x",
		"sid":  "$.cookies.sid",
		"uuid": "uuid()",
	}, Builtins{UUID: "u-1", Now: func() string { return "now" }})

	m, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("expected map, got %T", got)
	}
	if m["x"] != "1" {
		t.Fatalf("x: %#v", m["x"])
	}
	if m["sid"] != "abc" {
		t.Fatalf("sid: %#v", m["sid"])
	}
	if m["uuid"] != "u-1" {
		t.Fatalf("uuid: %#v", m["uuid"])
	}
}

func TestCaptureRequestBody_RestoresBody(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "http://example.com", strings.NewReader("hello"))
	_, b, truncated := CaptureRequestBody(r, 3)
	if string(b) != "hel" {
		t.Fatalf("b: %q", string(b))
	}
	if !truncated {
		t.Fatalf("expected truncated")
	}

	// Second read should see the same bytes we put back.
	_, b2, _ := CaptureRequestBody(r, 10)
	if string(b2) != "hel" {
		t.Fatalf("b2: %q", string(b2))
	}
}

func TestEvalJSONPathFromRequest_UsesCanonicalRoot(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "http://example.com/test?q=1&q=2", bytes.NewBufferString(`{"a":1}`))
	r.Header.Set("X-Test", "abc")
	r.AddCookie(&http.Cookie{Name: "sid", Value: "123"})

	if got := EvalJSONPathFromRequest(r, "$.headers.x-test", nil); got != "abc" {
		t.Fatalf("headers: %#v", got)
	}
	if got := EvalJSONPathFromRequest(r, "$.cookies.sid", nil); got != "123" {
		t.Fatalf("cookies: %#v", got)
	}
	if got := EvalJSONPathFromRequest(r, "$.query.q", nil); got != "1" {
		t.Fatalf("query first: %#v", got)
	}
	got := EvalJSONPathFromRequest(r, "$.query_all.q", nil)
	s, ok := got.([]string)
	if !ok || len(s) != 2 || s[0] != "1" || s[1] != "2" {
		t.Fatalf("query_all: %#v", got)
	}
	if got := EvalJSONPathFromRequest(r, "$.body.a", nil); got != float64(1) {
		t.Fatalf("body: %#v", got)
	}
}
