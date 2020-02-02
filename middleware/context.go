package middleware

import (
	"context"
	"github.com/foreversmart/plate/common"

	"github.com/gin-gonic/gin"
)

func SetRequestIdToContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString(common.RequestID)
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, common.RequestID, requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func RequestID(c *gin.Context) string {
	return c.GetString(common.RequestID)
}
