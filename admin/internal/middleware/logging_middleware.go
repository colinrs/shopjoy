package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

// LoggingMiddleware creates a middleware that logs HTTP requests and responses.
func LoggingMiddleware() rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Capture request body
			var reqBody []byte
			if r.Body != nil {
				reqBody, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			}

			// Create response capture
			respWriter := &responseWriter{
				ResponseWriter: w,
				body:           &bytes.Buffer{},
			}

			next(respWriter, r)

			duration := time.Since(start)
			reqBodyStr := ""
			if len(reqBody) > 0 {
				if len(reqBody) > 500 {
					reqBodyStr = string(reqBody[:500]) + "...(truncated)"
				} else {
					reqBodyStr = string(reqBody)
				}
			}
			respBodyStr := ""
			if respWriter.body.Len() > 0 {
				if respWriter.body.Len() > 500 {
					respBodyStr = respWriter.body.String()[:500] + "...(truncated)"
				} else {
					respBodyStr = respWriter.body.String()
				}
			}

			// Single line log
			logx.WithContext(r.Context()).Infof("%s %s | status=%d | duration=%v | req_body=%s | resp_body=%s",
				r.Method, r.URL.Path, respWriter.statusCode, duration, reqBodyStr, respBodyStr)
		}
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
