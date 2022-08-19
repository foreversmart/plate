package ginroute

import (
	"bytes"
	"fmt"
	"github.com/foreversmart/plate/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGinRecovery(t *testing.T) {
	type TestData struct {
		Int    int
		String string
	}

	b := bytes.NewBufferString("{a : \"a\"}")
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", b)

	args := []interface{}{&TestData{
		Int:    20,
		String: "first",
	},
		&TestData{
			Int:    11,
			String: "second",
		}}

	defer func() {
		if v := recover(); v != nil {
			_, err := GinRecovery(v, req, args)
			assert.NotNil(t, err)
			assert.Equal(t, errors.ErrorUnknownError, err)
		}
	}()

	panic(fmt.Errorf("panic error %s", "1"))
}
