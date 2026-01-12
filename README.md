# httprequestremap

Small Go helper package to:

- Capture an `*http.Request` body without consuming it.
- Build a canonical JSON object from an `*http.Request` (headers, query, cookies, body).
- Apply a structured remap template using JSONPath (same semantics used in network-observability).

## Install / Use

This repo uses multiple Go modules. For local development you can use a `replace`:

```go
require github.com/extedcouD/HttpRequestRemapper v0.0.0
```

## Example

```go
present, bodyBytes, _ := httprequestremap.CaptureRequestBody(r, 1024*1024)
_ = present

root := httprequestremap.RequestRoot(r, bodyBytes)

out := httprequestremap.ApplyTemplate(root, map[string]any{
  "sid": "$.cookies.sid",
  "ua": "$.headers.user-agent",
  "uuid": "uuid()",
}, httprequestremap.Builtins{UUID: "..."})
_ = out
```
