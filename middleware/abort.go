package middleware

import (
	"github.com/foreversmart/plate/controller/resp"
	"github.com/foreversmart/plate/errors"
	"github.com/gin-gonic/gin"
)

func abort(c *gin.Context, err *errors.Error) {
	c.AbortWithStatusJSON(err.Code, &resp.ResponseData{
		Error: err,
	})

}
