package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

func SetRequestIdToContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString(RequestIDKey)
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func RequestID(c *gin.Context) string {
	return c.GetString(RequestIDKey)
}
