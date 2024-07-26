# GoExpert Rate Limiter Challenge

This is a sample rate limiter implementation in Go.

## Example

The file `cmd/` has examples APIs using two Rate Limiters middlewares
with multiple settings.

- `cmd/memory/server.go`: uses in memory caching;
- `cmd/redis/server.go`: uses redis for caching;

The settings are defined with environment variables as in `.env` or `docker-compose.yml`.

```shell
# Lower priority Rate Limiter with IP as key (10 req/s).
IP_LIMIT=10
IP_LIMIT_DURATION=1 # in seconds

# Higher priority Rate Limiter with API-Key as key (100 req/s).
API_KEY_LIMIT=100
API_KEY_LIMIT_DURATION=1 # in seconds

# Redis config if running with Redis for caching
REDIS_ADDRESS=localhost:6379
REDIS_PASSWORD=
```

To run the example application you can run:
- Run with memory caching:
    - Setup te environment variables in `.env`
    - Run: `go run cmd/memory/server.go`
- Run with redis caching:
    - Setup te environment variables in `docker-compose.yml`
    - Run: `docker compose up -d`

Example requests are available at `api/requests.http`, or you can run the following commands:

```shell
# Lower priority Rate Limiter with IP as key.
curl -X GET http://localhost:8080/hello

# Higher priority Rate Limiter with API-Key as key
curl -X GET http://localhost:8080/hello -H "X-API-Key: 123456"
```

## Testing

To execute all the unit tests run 
```shell
go test ./... -v
```
