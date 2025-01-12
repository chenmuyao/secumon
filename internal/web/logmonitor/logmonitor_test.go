package logmonitor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chenmuyao/secumon/internal/domain"
	"github.com/chenmuyao/secumon/internal/event/monitor"
	monitormocks "github.com/chenmuyao/secumon/internal/event/monitor/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAccessLogAPI(t *testing.T) {
	testCases := []struct {
		Name string

		mock func(ctrl *gomock.Controller) monitor.LogMonitorPublisher

		// Inputs
		reqBuilder func(t *testing.T) *http.Request

		// Outputs
		wantCode int
		wantRes  Result
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
			mock: func(ctrl *gomock.Controller) monitor.LogMonitorPublisher {
				p := monitormocks.NewMockLogMonitorPublisher(ctrl)
				p.EXPECT().Publish(gomock.Any(), domain.AccessLog{
					Timestamp:  time.Date(2025, time.January, 8, 12, 0, 0, 0, time.UTC),
					ClientIP:   "192.168.1.1",
					Endpoint:   "/api/v1/resource",
					Method:     "GET",
					StatusCode: 401,
				}).Return(nil)
				return p
			},
			wantCode: http.StatusOK,
			wantRes:  ResultLogOK,
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
			mock: func(ctrl *gomock.Controller) monitor.LogMonitorPublisher {
				p := monitormocks.NewMockLogMonitorPublisher(ctrl)
				return p
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
			mock: func(ctrl *gomock.Controller) monitor.LogMonitorPublisher {
				p := monitormocks.NewMockLogMonitorPublisher(ctrl)
				p.EXPECT().Publish(gomock.Any(), domain.AccessLog{
					Timestamp:  time.Date(2025, time.January, 8, 12, 0, 0, 0, time.UTC),
					ClientIP:   "192.168.1.1",
					Endpoint:   "/api/v1/resource",
					Method:     "GET",
					StatusCode: 401,
				}).Return(nil)
				return p
			},
			wantCode: http.StatusOK,
			wantRes:  ResultLogOK,
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
			mock: func(ctrl *gomock.Controller) monitor.LogMonitorPublisher {
				p := monitormocks.NewMockLogMonitorPublisher(ctrl)
				return p
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
			mock: func(ctrl *gomock.Controller) monitor.LogMonitorPublisher {
				p := monitormocks.NewMockLogMonitorPublisher(ctrl)
				return p
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
			mock: func(ctrl *gomock.Controller) monitor.LogMonitorPublisher {
				p := monitormocks.NewMockLogMonitorPublisher(ctrl)
				return p
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			publisher := tc.mock(ctrl)
			// create the test server
			hdl := NewLogHandler(publisher, nil)
			server := gin.Default()

			hdl.RegisterHandlers(server)

			req := tc.reqBuilder(t)
			rec := httptest.NewRecorder()

			server.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantCode, rec.Code)
			var res Result
			if tc.wantRes != res {
				err := json.NewDecoder(rec.Body).Decode(&res)
				assert.NoError(t, err)
				assert.Equal(t, tc.wantRes, res)
			}
		})
	}
}

// limit unnecessary logs
func init() {
	gin.SetMode(gin.ReleaseMode)
}
