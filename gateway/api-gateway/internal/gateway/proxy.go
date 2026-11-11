package gateway

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func NewProxy(target string) gin.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Backend service unavailable"))
	}

	stripPrefix := http.StripPrefix("/api/v1", proxy)

	return func(c *gin.Context) {
		stripPrefix.ServeHTTP(c.Writer, c.Request)
	}
}
