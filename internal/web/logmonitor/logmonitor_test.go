package logmonitor

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAccessLogAPI(t *testing.T) {
	testCases := []struct {
		Name string

		// Inputs
		reqBuilder func(t *testing.T) *http.Request

		// Outputs
		wantCode int
	}{
		{
			Name: "test ok",
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/logs", bytes.NewBuffer([]byte(`{
"timestamp": "2025-01-08T12:00:00Z",
"client_ip": "192.168.1.1",
"endpoint": "/api/v1/resource",
"method": "GET",
"status_code": 401
}`)))
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// create the test server
			hdl := NewLogHandler()
			server := gin.Default()

			hdl.RegisterHandlers(server)

			req := tc.reqBuilder(t)
			rec := httptest.NewRecorder()

			server.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantCode, rec.Code)
		})
	}
}
