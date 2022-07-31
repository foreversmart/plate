package request

import (
	"github.com/valyala/fastjson"
	"io/ioutil"
	"net/http"
)

// Requester wrap http.Request and path param methods
type Requester interface {
	Request() *http.Request
	Params() map[string][]string // path param value
}

// fetchJsonAndMeta fetch json and meta from requester
func fetchJsonAndMeta(req Requester) (jsonValue *fastjson.Value, meta map[string]map[string][]string, err error) {
	body, err := ioutil.ReadAll(req.Request().Body)
	defer req.Request().Body.Close()
	if err != nil {
		return nil, nil, err
	}

	jsonValue, err = fastjson.Parse(string(body))
	if err != nil {
		return nil, nil, err
	}

	meta = make(map[string]map[string][]string)
	meta[LocHeader] = req.Request().Header

	formErr := req.Request().ParseForm()
	// indicate is form request
	if formErr == nil {
		meta[LocForm] = req.Request().PostForm
	}
	meta[LocQuery] = req.Request().URL.Query()
	meta[LocPath] = req.Params()

	return
}
