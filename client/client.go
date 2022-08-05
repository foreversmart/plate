package client

import (
	"bytes"
	"encoding/json"
	"github.com/foreversmart/plate/errors"
	"github.com/foreversmart/plate/logger"
	"io/ioutil"
	"net/http"
)

func Call(method, url string, req, resp interface{}, header http.Header, logger logger.Logger) (err error) {
	var request *http.Request

	if req != nil {
		bt, _ := json.Marshal(req)
		reader := bytes.NewReader(bt)
		request, _ = http.NewRequest(method, url, reader)
		request.Header.Set("Content-Type", "application/json")
	} else {
		request, _ = http.NewRequest(method, url, nil)
	}

	if header != nil {
		for k, v := range header {
			request.Header[k] = v
		}
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Error("post error", err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("read all error:", err)
		return err
	}

	var extErr errors.Error
	err = json.Unmarshal(body, &extErr)
	if err == nil && extErr.Code != 0 {
		logger.Error("service error:", extErr)
		return &extErr
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		logger.Error("unmarshal error:", err)
		return err
	}

	return
}
