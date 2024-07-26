package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mdalboni/goexpert-rate-limiter/internal/stores/memory"
	"github.com/mdalboni/goexpert-rate-limiter/pkg/limiter"
	"github.com/stretchr/testify/suite"
)

type RateLimiterMiddlewareTestSuite struct {
	suite.Suite
	handler    http.Handler
	middleware func(http.Handler) http.Handler
}

func (suite *RateLimiterMiddlewareTestSuite) SetupTest() {
	suite.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	store := memory.NewMemoryStore()

	limiter := limiter.NewLimiter(store, 1, time.Second)
	suite.middleware = NewRateLimiterMiddleware(limiter, nil)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(RateLimiterMiddlewareTestSuite))
}

func (suite *RateLimiterMiddlewareTestSuite) TestGivenMiddlewareWhenRequestIsExecutedBelowLimitShouldReturnStatusOk() {
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	req1.RemoteAddr = "192.168.0.1:12345"
	rec1 := httptest.NewRecorder()

	suite.middleware(suite.handler).ServeHTTP(rec1, req1)
	suite.Equal(http.StatusOK, rec1.Code)
}

func (suite *RateLimiterMiddlewareTestSuite) TestGivenMiddlewareWhenRequestIsExecutedAboveLimitShouldReturnStatusTooManyRequests() {
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	req1.RemoteAddr = "192.168.0.2:12345"

	// First request is allowed
	rec1 := httptest.NewRecorder()
	suite.middleware(suite.handler).ServeHTTP(rec1, req1)
	suite.Equal(http.StatusOK, rec1.Code)

	// Second request is limited
	rec2 := httptest.NewRecorder()
	suite.middleware(suite.handler).ServeHTTP(rec2, req1)
	suite.Equal(http.StatusTooManyRequests, rec2.Code)
}

func (suite *RateLimiterMiddlewareTestSuite) TestGivenMiddlewareWhenKeyIsEmptyThenShouldReturnStatusOk() {
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	req1.RemoteAddr = ""

	// First request is allowed
	rec1 := httptest.NewRecorder()
	suite.middleware(suite.handler).ServeHTTP(rec1, req1)
	suite.Equal(http.StatusOK, rec1.Code)

	// Second request should be limited but is allowed
	rec2 := httptest.NewRecorder()
	suite.middleware(suite.handler).ServeHTTP(rec2, req1)
	suite.Equal(http.StatusOK, rec2.Code)
}

func (suite *RateLimiterMiddlewareTestSuite) TestGivenMiddlewareWhenContextRateLimitAllowedKeyIsTrueThenShouldBypassLimiter() {
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	req1.RemoteAddr = "192.168.0.3:12345"

	// First request is allowed
	rec1 := httptest.NewRecorder()
	suite.middleware(suite.handler).ServeHTTP(rec1, req1)
	suite.Equal(http.StatusOK, rec1.Code)

	// Second request should be limited but is allowed with the context set
	ctx := req1.Context()
	ctx = context.WithValue(ctx, rateLimitAllowedKey, true)
	req2 := req1.WithContext(ctx)
	rec2 := httptest.NewRecorder()
	suite.middleware(suite.handler).ServeHTTP(rec2, req2)
	suite.Equal(http.StatusOK, rec2.Code)
}
