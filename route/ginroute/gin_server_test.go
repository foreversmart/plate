package ginroute

import (
	"github.com/foreversmart/plate/client"
	"github.com/foreversmart/plate/logger"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestNewGinServer(t *testing.T) {
	uid := "123456"

	server := NewGinServer()
	route := server.Route()
	route.Handle(http.MethodPost, "/price", _Business.Price, &PriceReq{})
	route.AddMiddleBefore(UserMiddleware, &UserReq{})
	go server.Run(":8080")

	time.Sleep(time.Second)

	req := &PriceReq{
		ProductID: 5,
	}
	var (
		resp *PriceResp
	)

	header := make(map[string][]string)
	header["uid"] = []string{uid}
	err := client.Call(http.MethodPost, "http://127.0.0.1:8080/price", req, &resp, header, logger.StdLog)
	assert.Nil(t, err)
	assert.Equal(t, req.ProductID*10, resp.Price)
	assert.Equal(t, "user:"+uid, resp.UserName)
}
