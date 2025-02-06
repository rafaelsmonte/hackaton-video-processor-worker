package httpServer

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartHTTPServer(t *testing.T) {
	os.Setenv("HTTP_PORT", "8084")
	defer os.Unsetenv("HTTP_PORT")

	req := httptest.NewRequest("GET", "/metrics", nil)
	rr := httptest.NewRecorder()

	go StartHTTPServer(nil)

	time.Sleep(time.Millisecond * 10)

	http.DefaultServeMux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.Equal(t, "OK", rr.Body.String())
}
