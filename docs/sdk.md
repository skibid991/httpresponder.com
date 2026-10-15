# httpresponder.com SDK reference

All clients use only the public HTTP API. No API key is needed.

| JavaScript | Python | Go | Behavior |
| --- | --- | --- | --- |
| `decode(data, contentType)` | `decode(data, content_type=...)` | `Decode(ctx, data, contentType)` | Returns bytes/Response; decode text explicitly if appropriate. |
| `encode(data)` | `encode(data)` | `Encode(ctx, data)` | Encode UTF-8 as base64. |
| `respond(options)` | `respond(**options)` | `Respond(ctx, params)` | GET custom response; nested headers are encoded as JSON. |
| `status(code)` | `status(code)` | `Status(ctx, code)` | Return a requested status without raising for HTTP error statuses. |
| `echo(options)` | `echo(**options)` | `Echo(ctx, Options)` | Echo the request; method/body/headers can be supplied. |
| `endpoint(path, options)` | `endpoint(path, **options)` | `Endpoint(ctx, path, Options)` | Access every other diagnostic route while preserving statuses and redirects. |

## Transport, errors and binary data

JavaScript uses native `Response`, including `.text()`, `.json()`, `.arrayBuffer()`, `.status` and `.headers`. Python returns `Response` with `.body` bytes, `.text()`, `.json()`, `.status` and `.headers`. Go returns `*Response` with `Body []byte`, `Text()`, `JSON(&value)`, `StatusCode` and `Header`.

Low-level request methods accept an API path, query parameters, method, headers and body. JavaScript/Python encode object bodies as JSON; Go uses `Options.Body` for JSON and `Options.RawBody` for bytes. Redirects remain visible by default. Non-2xx statuses raise `APIError` unless disabled through JavaScript/Python `check: false` / `check=False` or Go `AllowHTTPError: true`. Methods intended to return arbitrary HTTP statuses disable that check automatically. Network errors and timeouts remain ordinary runtime errors.

Clients do not perform automatic retries. A timeout does not prove a creation request failed. Binary responses must be saved as bytes. These convenience clients buffer responses; use your language's native HTTP streaming client for long-lived SSE or incremental reads.

## Custom settings

JavaScript: `new HTTPResponder({timeout: 60000})`; inject a fetch implementation with `fetch: myFetch`.

Python: `HTTPResponder(timeout=60)`.

Go: `client := httpresponder.New()`; configure `client.HTTPClient.Timeout` and pass a cancellable context to each method. Timeout units are milliseconds in JavaScript, seconds in Python, and `time.Duration` in Go.

For testing, override JavaScript `baseUrl`, Python `base_url`, or Go `BaseURL`. Production defaults to `https://httpresponder.com`.
