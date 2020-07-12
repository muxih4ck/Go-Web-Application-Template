package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"go.uber.org/zap"

	"github.com/muxih4ck/Go-Web-Application-Template/handler"
	"github.com/muxih4ck/Go-Web-Application-Template/log"
	"github.com/muxih4ck/Go-Web-Application-Template/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/willf/pad"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logging is a middleware function that logs the each request.
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UTC()
		path := c.Request.URL.Path

		// Skip for the health check requests.
		if path == "/sd/health" || path == "/sd/ram" || path == "/sd/cpu" || path == "/sd/disk" {
			return
		}

		// Read the Body content
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}

		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// The basic informations.
		method := c.Request.Method
		ip := c.ClientIP()

		//log.Debugf("New request come in, path: %s, Method: %s, body `%s`", path, method, string(bodyBytes))
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		// Continue.
		c.Next()

		// Calculates the latency.
		end := time.Now().UTC()
		latency := end.Sub(start)

		code, message := -1, ""

		// get code and message
		var response handler.Response
		if err := json.Unmarshal(blw.body.Bytes(), &response); err != nil {
			log.Error("JSON unmarshal failed",
				zap.String("reason", err.Error()),
				zap.String("detail", fmt.Sprintf("response body can not unmarshal to model.Response struct, body: `%s`", blw.body.Bytes())))
			code = errno.InternalServerError.Code
			message = err.Error()
		} else {
			code = response.Code
			message = response.Message
		}

		log.Info(
			fmt.Sprintf("%-13s | %-12s | %s %s | {code: %d, message: %s}",
				latency,
				ip,
				pad.Right(method, 5, ""),
				path,
				code,
				message))
	}
}
