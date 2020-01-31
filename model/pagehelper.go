package model

import "github.com/foreversmart/plate/model/errors"

// 页面参数
const (
	MinPageSize     = 10
	MaxPageSize     = 100
	DefaultPageSize = 20
)

// PagingHelper 分页辅助
func PagingHelper(pageNum, pageSize int) (offset int, err error) {
	if pageNum <= 0 {
		return 0, errors.ErrInvalidParams
	}
	if pageSize < MinPageSize {
		pageSize = MinPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	offset = 0
	offset = (pageNum - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	return
}
