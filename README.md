# HTTP testing API: base64, echo and custom responses

[HTTPResponder](https://httpresponder.com) is a free HTTP request and response testing service. Decode base64 from a URL, inspect request headers and bodies, return custom status codes, or test redirects and streaming. No account or API key is required. Most diagnostics work with a single GET request.

## Decode base64 from a URL

```sh
curl 'https://httpresponder.com/base64/SGVsbG8sIHdvcmxkIQ=='
# Hello, world!

curl --get 'https://httpresponder.com/base64' \
  --data-urlencode 'data=Hello 🌎' --data-urlencode 'mode=encode'
```

Both standard and URL-safe base64 are accepted, padded or unpadded. Prefer URL-safe encoding in paths. For query parameters, use a proper query encoder so `+` does not become a space. Add `content_type=application/octet-stream` when decoding binary content; preserve bytes rather than decoding them as text. `/base64` also accepts POST JSON with `data`, `mode` and `content_type`.

## Return a custom response

```sh
curl -i --get 'https://httpresponder.com/respond' \
  --data-urlencode 'status=201' --data-urlencode 'body=Created' \
  --data-urlencode 'headers={"X-Example":"yes"}'

curl -i 'https://httpresponder.com/mix/s=201/h=content-type:text%2Fplain/b64=SGVsbG8='
curl -i 'https://httpresponder.com/status/418'
```

`/respond` accepts GET query parameters or POST JSON: `status`, `body`, `base64` (takes precedence over body), `content_type`, `headers` and `delay` in seconds. Status codes range from 200 to 599; statuses without bodies retain HTTP semantics. Delay is at most 10 seconds. Runtime-managed headers cannot be overridden. HTML responses may run JavaScript when opened in a browser; fetching text alone does not execute it.

`/mix/` combines slash-separated directives: `s` for status, `h` for `name:value` headers, `c` for cookies, `r` for redirect URL, `b64` for body, and `end` to stop parsing. URL-encode slashes inside values. Go-style templates are not supported.

## Echo a request

```sh
curl 'https://httpresponder.com/anything?example=yes'
curl 'https://httpresponder.com/anything' -H 'Content-Type: application/json' -d '{"hello":"world"}'
```

The response includes `method`, `url`, `headers`, `args`, raw `data`, parsed `json`, `form`, `files` and `origin`. `/anything` accepts different HTTP methods. `/get`, `/post`, `/put`, `/patch` and `/delete` test their named methods. Avoid submitting real credentials or confidential data to a public diagnostic service.

## Endpoint directory

| Endpoint | Purpose |
| --- | --- |
| `/headers`, `/ip`, `/user-agent`, `/uuid` | Inspect headers, address, user agent; generate a UUID |
| `/status/200,503` | Randomly select a listed status |
| `/delay/1` | Delay up to 10 seconds, then echo |
| `/bytes/32` | Random bytes, up to 1 MiB |
| `/stream/5` | NDJSON, up to 100 records |
| `/sse?count=5&duration=2` | Server-sent events |
| `/drip?numbytes=10&duration=2` | Drip response bytes; total delay plus duration at most 10 seconds |
| `/redirect/3`, `/relative-redirect/3`, `/absolute-redirect/3` | Redirect chains, at most 20 hops |
| `/redirect-to?url=...&status_code=302` | Redirect to a selected destination |
| `/response-headers?X-Example=yes` | Return query parameters as response headers |
| `/etag/example`, `/cache/60` | Conditional requests and cache headers |
| `/basic-auth/user/password`, `/bearer/example` | Test Basic or Bearer authentication |
| `/cookies`, `/cookies/set?example=yes`, `/cookies/delete?example=` | Inspect, set or remove cookies |

## Common questions

**Does a 418 response mean the API failed?** When testing a status, it is the requested result. SDK `status`, `respond` and `endpoint` methods preserve non-2xx responses instead of throwing, and leave redirects visible. Read `.status` in JavaScript/Python or `.StatusCode` in Go.

**Can a GET-only client use this?** Yes: base64 decoding, custom bodies, status codes, headers and most diagnostics have URL-based forms. Echoing a POST body naturally requires POST.

## SDKs and runnable examples

This repository contains dependency-free JavaScript (with TypeScript declarations), Python and Go clients. They are source packages in this repository; no npm or PyPI release is implied. Clone the repository to run the examples:

```sh
git clone https://github.com/skibid991/httpresponder.com.git
cd httpresponder.com
```

### JavaScript and TypeScript

Requires Node.js 22+ or a browser with `fetch` and `AbortSignal.timeout`. Run `node examples/javascript/quickstart.mjs` from the repository root. You can also install the local package using `npm install ./sdk/javascript`; TypeScript declarations are included.

```js
import {HTTPResponder} from './sdk/javascript/index.js';
const client = new HTTPResponder();
console.log(await (await client.decode('SGVsbG8sIHdvcmxkIQ==')).text());
const response = await client.status(418);
console.log('Requested status:', response.status);
```

### Python

Requires Python 3.10+, with no third-party packages. Run `PYTHONPATH=sdk/python python3 examples/python/quickstart.py`. To use the module in another project, copy `sdk/python/httpresponder.py` into your project or add that directory to your Python path.

```python
from httpresponder import HTTPResponder

client = HTTPResponder()
print(client.decode('SGVsbG8sIHdvcmxkIQ==').text())
print('Requested status:', client.status(418).status)
```

### Go

Requires Go 1.22+. Run `go -C examples/go run .`. The example module uses the local SDK. From another Go module, run `go get github.com/skibid991/httpresponder.com/sdk/go` and import that path.

See [the runnable Go example](examples/go/main.go) and [SDK method reference](docs/sdk.md). SDKs use a 45-second client timeout by default. Change it for your workload. Do not automatically retry creation requests: a successful request whose response is lost can still have created a resource.

## API reference and discovery

- [Interactive HTML API reference](https://httpresponder.com/docs)
- [OpenAPI 3.1 schema](https://httpresponder.com/openapi.json)
- [Plain-text usage guide](https://httpresponder.com/llms.txt)

The live API reference is authoritative. Documented GET routes work without an SDK. Encode query parameters with `URLSearchParams`, Python `urlencode`, Go `url.Values`, or cURL `--data-urlencode`.

## Development

Offline tests cover query encoding, Unicode, binary response preservation, request bodies, authorization headers, error handling and redirect behavior:

```sh
node --test sdk/javascript/client.test.mjs
PYTHONPATH=sdk/python python3 -m unittest discover -s sdk/python -p 'test_*.py'
go -C sdk/go test ./...
go -C examples/go build -o /tmp/httpresponder-example .
```

Examples call the public service. Paste/link quickstarts create temporary resources and delete them after reading. Contributions should focus on public API usage and client behavior. Licensed under [MIT](LICENSE).
