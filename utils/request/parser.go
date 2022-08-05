package request

import (
	"github.com/foreversmart/plate/utils/view"
	"github.com/valyala/fastjson"
	"reflect"
)

type Parser struct {
	jsonValue *fastjson.Value
	meta      map[string]map[string][]string
	mid       *view.View
}

func NewParser(req Requester) (p *Parser, err error) {
	p = &Parser{}
	p.jsonValue, p.meta, err = fetchJsonAndMeta(req)
	return
}

func (p *Parser) WithMid(resp interface{}) error {
	v, err := view.FetchViewFromStruct(reflect.ValueOf(resp), TagNameFetch, LocMid)
	if err != nil {
		return err
	}

	if p.mid == nil {
		p.mid = v
		return nil
	}

	p.mid.MergeWithNew(v)

	return nil
}

func (p *Parser) Parse(v reflect.Value) error {
	return Parse(v, p.jsonValue, p.meta, p.mid, nil)
}
