package gin_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	libgin "github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/middleware/gin"
	"github.com/ulule/limiter/drivers/store/memory"
)

func TestHTTPMiddleware(t *testing.T) {
	is := require.New(t)
	libgin.SetMode(libgin.TestMode)

	request, err := http.NewRequest("GET", "/", nil)
	is.NoError(err)
	is.NotNil(request)

	store := memory.NewStore()
	is.NotZero(store)

	rate, err := limiter.NewRateFromFormatted("10-M")
	is.NoError(err)
	is.NotZero(rate)

	middleware := gin.NewMiddleware(limiter.New(store, rate))
	is.NotZero(middleware)

	router := libgin.New()
	router.Use(middleware)
	router.GET("/", func(c *libgin.Context) {
		c.String(http.StatusOK, "hello")
	})

	success := int64(10)
	clients := int64(100)

	//
	// Sequential
	//

	for i := int64(1); i <= clients; i++ {

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, request)

		if i <= success {
			is.Equal(resp.Code, http.StatusOK)
		} else {
			is.Equal(resp.Code, http.StatusTooManyRequests)
		}
	}

	//
	// Concurrent
	//

	store = memory.NewStore()
	is.NotZero(store)

	middleware = gin.NewMiddleware(limiter.New(store, rate))
	is.NotZero(middleware)

	router = libgin.New()
	router.Use(middleware)
	router.GET("/", func(c *libgin.Context) {
		c.String(http.StatusOK, "hello")
	})

	wg := &sync.WaitGroup{}
	counter := int64(0)

	for i := int64(1); i <= clients; i++ {
		wg.Add(1)
		go func() {

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, request)

			if resp.Code == http.StatusOK {
				atomic.AddInt64(&counter, 1)
			}

			wg.Done()
		}()
	}

	wg.Wait()
	is.Equal(success, atomic.LoadInt64(&counter))

	//
	// Custom KeyGetter
	//

	store = memory.NewStore()
	is.NotZero(store)

	j := 0
	KeyGetter := func(c *libgin.Context) string {
		j += 1
		return strconv.Itoa(j)
	}
	middleware = gin.NewMiddleware(limiter.New(store, rate), gin.WithKeyGetter(KeyGetter))

	is.NotZero(middleware)

	router = libgin.New()
	router.Use(middleware)
	router.GET("/", func(c *libgin.Context) {
		c.String(http.StatusOK, "hello")
	})

	for i := int64(1); i <= clients; i++ {
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, request)
		// We should always be ok as the key changes for each request
		is.Equal(http.StatusOK, resp.Code, strconv.Itoa(int(i)))
	}

}
