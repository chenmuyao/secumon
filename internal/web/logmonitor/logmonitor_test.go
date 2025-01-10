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
		wantRes  string
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
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantRes: `{
"status": "success",
"message": "Log received and queued."
}`,
		},
		{
			Name: "wrong json format",
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/logs", bytes.NewBuffer([]byte(`{
"timestamp": "2025-01-08T12:00:00Z",
"client_ip": "192.168.1.1",
"endpoint": "/api/v1/resource",
"method": "GET",
"status_code401
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusBadRequest,
		},
		{
			Name: "unknown fields, just ignore but ok",
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/logs", bytes.NewBuffer([]byte(`{
"timestamp": "2025-01-08T12:00:00Z",
"client_ip": "192.168.1.1",
"endpoint": "/api/v1/resource",
"method": "GET",
"status_code": 401,
"other_things": "something"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantRes: `{
"status": "success",
"message": "Log received and queued."
}`,
		},
		{
			Name: "missing fields",
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/logs", bytes.NewBuffer([]byte(`{
"client_ip": "192.168.1.1",
"endpoint": "/api/v1/resource",
"method": "GET",
"status_code": 401
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusBadRequest,
		},
		{
			Name: "wrong time format",
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/logs", bytes.NewBuffer([]byte(`{
"timestamp": "2025-01-08",
"client_ip": "192.168.1.1",
"endpoint": "/api/v1/resource",
"method": "GET",
"status_code": 401
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusBadRequest,
		},
		{
			Name: "wrong client IP",
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/logs", bytes.NewBuffer([]byte(`{
"timestamp": "2025-01-08T12:00:00Z",
"client_ip": "192.168.1",
"endpoint": "/api/v1/resource",
"method": "GET",
"status_code": 401
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusBadRequest,
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

// limit unnecessary logs
func init() {
	gin.SetMode(gin.ReleaseMode)
}
