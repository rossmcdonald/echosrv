# JSON Echo

Just what it sounds like, a small, performant JSON-formatted echo service. This service can be used for debugging or demonstration purposes when inspecting the request data is important (or you simply don't care about the back-end logic). The container image also clocks in at about ~20MB in size, making it faster to deploy and in a smaller footprint than some of the alternatives out there.

When receiving an HTTP request, it will return and log a JSON document with the details of the request and some other pieces of metadata that are useful (hostname, timestamp, and parse/request errors). For example:

```json
% curl localhost:8889
{
  "host": "echo.internal",
  "ts": "2021-03-23T15:17:18.439079Z",
  "request": {
    "headers": {
      "Accept": [
        "*/*"
      ],
      "Content-Length": [
        "29"
      ],
      "Content-Type": [
        "application/json"
      ],
      "User-Agent": [
        "insomnia/2021.1.0"
      ]
    },
    "url": {
      "Scheme": "",
      "Opaque": "",
      "User": null,
      "Host": "",
      "Path": "/something",
      "RawPath": "",
      "ForceQuery": false,
      "RawQuery": "true=false",
      "Fragment": "",
      "RawFragment": ""
    },
    "body": {
      "get some?": "get some!"
    },
    "host": "localhost:8889",
    "proto": "HTTP/1.1",
    "method": "POST",
    "form": {
      "true": [
        "false"
      ]
    }
  },
  "errors": []
}
```

The same is logged server-side, allowing you to inspect request payloads whenever you need to (for example testing a webhook or collecting a header value).

## Running

See the `Makefile` for more details on how to build from source (you'll need Go). To run in Docker:

```sh
docker run -p 8889:8889 rossmcd/echo-json
```

There is also a sample `k8s.yaml` file as well for a generic Kubernetes service:

```sh
kubectl apply -f k8s.yaml
```

Once setup, run `curl localhost:8889` to generate a sample response.

## Metrics

Prometheus metrics are also included! Issue a GET request `/metrics` to get some metrics for a dashboard.

## Reference

See the `openapi.yaml` for an OpenAPI specification of the response format.
