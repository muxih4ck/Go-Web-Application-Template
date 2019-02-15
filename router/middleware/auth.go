package middleware

import (
	"github.com/muxih4ck/Go-Web-Application-Template/handler"
	"github.com/muxih4ck/Go-Web-Application-Template/pkg/errno"
	"github.com/muxih4ck/Go-Web-Application-Template/pkg/token"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the json web token.
		if _, err := token.ParseRequest(c); err != nil {
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
