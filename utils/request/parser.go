package request

import (
	"encoding/json"
	"github.com/valyala/fastjson"
	"reflect"
)

type Parser struct {
	jsonValue *fastjson.Value
	meta      map[string]map[string][]string
	mid       map[string]*fastjson.Value
}

func NewParser(req Requester) (p *Parser, err error) {
	p = &Parser{}
	p.jsonValue, p.meta, err = fetchJsonAndMeta(req)
	p.mid = make(map[string]*fastjson.Value)
	return
}

func (p *Parser) WithMid(resp interface{}) error {
	t := reflect.TypeOf(resp).String()
	jsonStr, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	jsonValue, err := fastjson.ParseBytes(jsonStr)
	if err != nil {
		return err
	}

	p.mid[t] = jsonValue
	return nil
}

func (p *Parser) Parse(v reflect.Value) error {
	return Parse(v, p.jsonValue, p.meta, nil)
}
