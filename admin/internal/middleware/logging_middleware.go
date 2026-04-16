package middleware

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

const (
	maxBodySize = 64 * 1024 // 64KB max for request/response body logging
	logTruncate = 500       // 500 bytes truncation for logged bodies
)

// sensitiveFields 敏感字段列表，用于日志脱敏
var sensitiveFields = []*regexp.Regexp{
	regexp.MustCompile(`"password"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`"token"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`"access_token"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`"refresh_token"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`"secret"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`"api_key"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`"authorization"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`"credit_card"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`"card_number"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`"cvv"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`"pin"\s*:\s*"[^"]*"`),
}

// sanitizeForLogging 对日志内容进行敏感信息脱敏
func sanitizeForLogging(data []byte) []byte {
	result := make([]byte, len(data))
	copy(result, data)
	for _, pattern := range sensitiveFields {
		result = pattern.ReplaceAll(result, []byte(`"[REDACTED]"`))
	}
	return result
}

// truncateBody 截断过长的body用于日志
func truncateBody(data []byte, maxLen int) string {
	if len(data) > maxLen {
		return string(data[:maxLen]) + "...(truncated)"
	}
	return string(data)
}

// LoggingMiddleware creates a middleware that logs HTTP requests and responses.
func LoggingMiddleware() rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Capture request body with size limit
			var reqBody []byte
			if r.Body != nil {
				// Limit body size to prevent memory exhaustion
				limitedReader := io.LimitReader(r.Body, maxBodySize)
				reqBody, _ = io.ReadAll(limitedReader)
				r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			}

			// Create response capture with size limit
			respWriter := &responseWriter{
				ResponseWriter: w,
				body:           &bytes.Buffer{},
				maxSize:        maxBodySize,
			}

			next(respWriter, r)

			duration := time.Since(start)

			// Sanitize and truncate for logging
			sanitizedReq := sanitizeForLogging(reqBody)
			reqBodyStr := truncateBody(sanitizedReq, logTruncate)

			sanitizedResp := sanitizeForLogging(respWriter.body.Bytes())
			respBodyStr := truncateBody(sanitizedResp, logTruncate)

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
	maxSize    int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	// Limit response body buffer to prevent memory exhaustion
	if rw.body.Len() < rw.maxSize {
		remaining := rw.maxSize - rw.body.Len()
		if len(b) > remaining {
			rw.body.Write(b[:remaining])
		} else {
			rw.body.Write(b)
		}
	}
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
