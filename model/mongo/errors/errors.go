package errors

import (
	"errors"

	"github.com/foreversmart/mgo"
)

var (
	ErrInvalidId     = errors.New("Invalid BSON ID")
	ErrInvalidParams = errors.New("Invalid params")
	ErrNotPersisted  = errors.New("Not persisted")
	ErrNotFound      = mgo.ErrNotFound
)
