# api-mock

A simple API mock only with Go and its standard library. 😎

Respond with **pre-defined JSON** to HTTP GET requests for **pre-defined entry points**.

# Getting Started 🎉

```zsh
$ docker compose up # or `go run main.go'
$ curl -i http://localhost:8080/v1/example?a=AAA\&b=BBB\&c=CCC
```

Take a look at the response and API logs. 👀

# Configure

### `.env`

- `API_ENTRY_POINT`: Set the _entry point_ for the API. Start with `/`.
- `API_PORT`: Set the _port_ for the API. Integer value.
- `TIME_ZONE`: Set the _time zone_ for the API. IANA time zone identifier. Mainly used for log output.

### `response.json`

- `status`: Set the _status_ of the response. Integer value.
- `header`: Set the _header_ of the response. key-value format for string only.
  - ⚠️ For `Content-Type`, practically only `application/json` and `text/plain` are supported.
- `body`: Set the _body_ of the response. String or JSON format.
  - ⚠️ For string, please set `Content-Type` to `text/plain`.
  - ⚠️ For JSON format, please set `Content-Type` to `application/json`.
- `no_content`: Set the presence or absence of the response body. `true` or `false`.

# License

MIT
