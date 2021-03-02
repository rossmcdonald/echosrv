# JSON Echo

Just what it sounds like, a JSON-formatted echo service. This service can be used for debugging or demonstration purposes when inspecting the request data is important (or you simply don't care about the back-end logic). 

When receiving an HTTP request, it will return and log a JSON document with the details of the request. For example:

```json
% curl localhost:8888
{
        "ts": "2021-03-02T17:54:46.622134Z",
        "request": {
                "headers": {
                        "Accept": [
                                "*/*"
                        ],
                        "User-Agent": [
                                "curl/7.64.1"
                        ]
                },
                "url": {
                        "Scheme": "",
                        "Opaque": "",
                        "User": null,
                        "Host": "",
                        "Path": "/",
                        "RawPath": "",
                        "ForceQuery": false,
                        "RawQuery": "",
                        "Fragment": "",
                        "RawFragment": ""
                },
                "body": {},
                "host": "localhost:8888",
                "proto": "HTTP/1.1",
                "method": "GET",
                "form": {}
        }
}
```

The same is logged server-side, allowing you to inspect request payloads whenever you need to (for example testing a webhook).

## Running

See the `Makefile` for more details on how to build from source (you'll need Go). To run in Docker:

```sh
docker run -p 8888:8888 rossmcd/echo-json
```

There is also a sample `k8s.yaml` file as well for a generic Kubernetes service:

```sh
kubectl apply -f k8s.yaml
```

Once setup, run `curl localhost:8888` to generate a sample response.

## Reference

See the `openapi.yaml` for an OpenAPI specification of the response format.
