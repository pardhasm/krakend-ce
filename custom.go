package krakend

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type customLog struct {
	Timestamp    time.Time     `json:"@timestamp"`
	Version      int           `json:"@version"`
	Level        string        `json:"level"`
	Client       string        `json:"client"`
	Method       string        `json:"method"`
	Path         string        `json:"path"`
	Proto        string        `json:"proto"`
	StatusCode   int           `json:"status"`
	Latency      time.Duration `json:"latency"`
	UserAgent    string        `json:"agent"`
	ErrorMessage string        `json:"error"`
	Header       http.Header   `json:"headers"`
}

func formatLogStashMessage(clientIP string, timeStamp time.Time, method string, path string, request *http.Request, statusCode int, latency time.Duration, userAgent string, errorMessage string) string {
	innerLog := &customLog{
		Client:       clientIP,
		Timestamp:    timeStamp,
		Version:      1,
		Level:        "INFO",
		Method:       method,
		Path:         path,
		Proto:        request.Proto,
		StatusCode:   statusCode,
		Latency:      latency,
		UserAgent:    userAgent,
		ErrorMessage: errorMessage,
		Header:       request.Header,
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(request.Body)
	reqStr := buf.String()

	jinnerlog, err := json.Marshal(innerLog)
	outerLog := &LogstashPattern{
		Timestamp:    timeStamp,
		Version:      1,
		Level:        "INFO",
		Message:      reqStr,
		Client:       clientIP,
		Module:       "[KRAKEND]",
		Method:       method,
		Path:         path,
		Proto:        request.Proto,
		StatusCode:   statusCode,
		Latency:      latency,
		UserAgent:    userAgent,
		ErrorMessage: errorMessage,
		APIKey:       request.Header.Get("Api-Key"),
	}

	jsonOuterLog, err := json.Marshal(outerLog)
	if err != nil {
		return "<Logging Error. Unable to log. Check log config>"
	}

	return strings.Replace(string(jsonOuterLog), "kkk", string(jinnerlog), 1)
}

func formatLog() func(param gin.LogFormatterParams) string {
	return func(param gin.LogFormatterParams) string {
		a := formatLogStashMessage(
			param.ClientIP,
			param.TimeStamp,
			param.Method,
			param.Path,
			param.Request,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage) + "\n"
		return a
	}
}

type LogstashPattern struct {
	Timestamp    time.Time     `json:"@timestamp"`
	Version      int           `json:"@version"`
	Level        string        `json:"level"`
	Message      string        `json:"message"`
	Module       string        `json:"module"`
	Client       string        `json:"client"`
	Method       string        `json:"method"`
	Path         string        `json:"path"`
	Proto        string        `json:"proto"`
	StatusCode   int           `json:"status"`
	Latency      time.Duration `json:"latency"`
	UserAgent    string        `json:"agent"`
	ErrorMessage string        `json:"error"`
	APIKey       string        `json:"api-key"`
}
