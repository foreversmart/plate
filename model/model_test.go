package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModel_Copy(t *testing.T) {
	var err error
	session1 := mongo.Session()
	session2 := session1.Copy()

	err = session1.Ping()

	assert.Nil(t, err)

	err = session2.Ping()
	assert.Nil(t, err)

	session2.Close()

	err = session1.Ping()
	assert.Nil(t, err)
}

func TestModel_C(t *testing.T) {

	var err error
	copiedDB := mongo.C(mongo.Config().Database)

	err = copiedDB.Session().Ping()
	assert.Nil(t, err)
	copiedDB.Close()

	err = mongo.Session().Ping()
	assert.Nil(t, err)
}
