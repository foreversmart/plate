package resp

import "github.com/foreversmart/plate/errors"

type ResponseData struct {
	Data      interface{}   `json:"data,omitempty"`
	Page      *int          `json:"page,omitempty"`
	Size      *int          `json:"size,omitempty"`
	Total     *int          `json:"total,omitempty"`
	Error     *errors.Error `json:"error,omitempty"`
	RequestId string        `json:"request_id,omitempty"`
}
