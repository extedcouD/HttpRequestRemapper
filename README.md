# HttpRequestRemapper

Small Go helper package to:

- Capture an `*http.Request` body without consuming it.
- Build a canonical JSON object from an `*http.Request` (headers, query, cookies, body).
- Evaluate JSONPath or apply a structured remap template directly against an `*http.Request`.

## Install / Use

This repo uses multiple Go modules. For local development you can use a `replace`:

```go
require github.com/extedcouD/HttpRequestRemapper v0.0.0
```

## Example

```go
// Evaluate a JSONPath directly on the request.
sid := httprequestremap.EvalJSONPathFromRequest(r, "$.cookies.sid", nil)
_ = sid

// Apply a template directly on the request.
out := httprequestremap.ApplyTemplateFromRequest(r, map[string]any{
  "sid":  "$.cookies.sid",
  "ua":   "$.headers.user-agent",
  "uuid": "uuid()",
}, httprequestremap.Builtins{UUID: "..."}, nil)
_ = out
```
