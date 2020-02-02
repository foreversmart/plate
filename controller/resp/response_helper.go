package resp

import (
	"github.com/foreversmart/plate/common"
	"github.com/foreversmart/plate/errors"
	"github.com/foreversmart/plate/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response 正常返回，提供给外部调用
func Response(c *gin.Context, data interface{}) {
	requestID := c.GetString(common.RequestID)
	CommonSetHeader(c)
	c.JSON(http.StatusOK, &ResponseData{
		Data:      data,
		RequestId: requestID,
	})
}

// PageResponse 分页返回，提供给外部调用
func PageResponse(c *gin.Context, data interface{}, page, size, total int) {
	requestID := c.GetString(common.RequestID)
	CommonSetHeader(c)
	c.JSON(http.StatusOK, &ResponseData{
		Data:      data,
		Page:      &page,
		Size:      &size,
		Total:     &total,
		RequestId: requestID,
	})
}

// ErrorResponse 错误返回，提供给外部调用
func ErrorResponse(c *gin.Context, err error) {
	requestID := c.GetString(common.RequestID)
	CommonSetHeader(c)
	logger.LoggerFromContext(c).Infof("bad response %s request id %s", err.Error(), requestID)
	if stdErr, ok := err.(*errors.Error); ok {
		c.JSON(stdErr.Code, &ResponseData{
			Error:     stdErr,
			RequestId: requestID,
		})

		return
	}

	logger.LoggerFromContext(c).Warnf("cant recognise warning %s", err)

	c.JSON(errors.InternalError.Code, &ResponseData{
		Error:     errors.InternalError,
		RequestId: requestID,
	})
}

func RawJSONResponse(c *gin.Context, data interface{}) {
	CommonSetHeader(c)
	c.JSON(http.StatusOK, data)
}

func RawCSVResponse(c *gin.Context, data []byte) {
	CommonSetHeader(c)
	c.Data(http.StatusOK, "text/csv", data)
}

func CommonSetHeader(c *gin.Context) {
	requestId := c.GetString(common.RequestID)
	c.Writer.Header().Set(common.RequestID, requestId)

	// TODO uniform time key
	c.Writer.Header().Set("data", time.Now().Format("2006-01-02T15:04:05.000+08:00"))
}
