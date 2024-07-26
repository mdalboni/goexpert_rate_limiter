package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mdalboni/goexpert-rate-limiter/internal/stores/redis"
	"github.com/mdalboni/goexpert-rate-limiter/pkg/config"
	"github.com/mdalboni/goexpert-rate-limiter/pkg/limiter"
	"github.com/mdalboni/goexpert-rate-limiter/pkg/middlewares"
)

func main() {
	cfg := config.LoadConfig()
	log.Println("Config loaded:", cfg)
	store := redis.NewRedisStore(cfg.RedisAddress, cfg.RedisPassword)

	ipLimiter := limiter.NewLimiter(
		store,
		cfg.IPLimit,
		cfg.IPDuration,
	)
	apiKeyLimiter := limiter.NewLimiter(
		store,
		cfg.APIKeyLimit,
		cfg.APIKeyDuration,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", helloRoute)

	ipRateLimiterMiddleware := middlewares.NewRateLimiterMiddleware(ipLimiter, nil)
	apiKeyRateLimiterMiddleware := middlewares.NewRateLimiterMiddleware(apiKeyLimiter, keyMapper)
	mid := ipRateLimiterMiddleware(mux)
	mid = apiKeyRateLimiterMiddleware(mid)
	mid = middlewares.LogRequest(mid)
	log.Println("Server is starting on port 8080")
	err := http.ListenAndServe(":8080", mid)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	log.Println("Server stopped")
}

func helloRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Welcome to the Rate Limiter API!"}`))
}

func keyMapper(r *http.Request) string {
	apiKey := r.Header.Get("API_KEY")
	if apiKey != "" {
		return fmt.Sprintf("API_KEY:%v", apiKey)
	}
	return ""
}
