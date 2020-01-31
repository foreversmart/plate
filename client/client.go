package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"server/errors"
	"server/utils"
)

func DoRequest(logger utils.Logger, uri string, args, ret interface{}) error {
	bt, _ := json.Marshal(args)
	reader := bytes.NewReader(bt)
	resp, err := http.Post(uri, "application/json", reader)
	if err != nil {
		logger.Error("post error", err)
		return err
	}
	defer resp.Body.Close()
	bt, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("read all error:", err)
		return err
	}
	var extErr errors.Error
	err = json.Unmarshal(bt, &extErr)
	if err == nil && extErr.Code != 0 {
		logger.Error("service error:", extErr)
		return &extErr
	}
	err = json.Unmarshal(bt, ret)
	if err != nil {
		logger.Error("unmarshal error:", err)
		return err
	}
	return nil
}
