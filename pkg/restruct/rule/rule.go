package rule

import (
	"errors"

	"github.com/tidwall/gjson"
)

// define field type, used to simulate constants
type FieldType string

const (
	NullType   FieldType = "null"
	StringType FieldType = "string"
	NumberType FieldType = "number"
	BoolType   FieldType = "bool"
	JsonType   FieldType = "json"
)

/*
Node of extract structure
Notion: Target, Extract, Fields, only one of them should appeared
*/
type Field struct {
	// field key, this should be unique in current level
	Key string
	// value type, only few types are supported
	Type FieldType
	// this mark final mapping field, if appeares, will stopped and return
	Target string
	// continue extract info from this field
	Extract Extract
	// children of node, if it's json type, maybe null
	Fields []Field
}

type Extract struct {
	// a function name to handle current field
	Action string
	// function params, only string array allowed
	Params []string
	// mark function return type, only few types allowed
	Return FieldType
	// this mark final mapping field, if appeares, will stopped and return
	Target string
	// mark result struct, it's array object, even only one field
	Fields []Field
}

func Parse(s string) (Extract, error) {
	return parseExtract(gjson.Parse(s))
}

func parseExtract(gr gjson.Result) (Extract, error) {
	var r gjson.Result
	var e Extract

	// this must a field struct, or throw error
	if !gr.IsObject() {
		return e, errors.New("不是一个对象，不能解析")
	}

	r = gr.Get("action")
	if r.Type != gjson.String {
		return e, errors.New("action 字段必须是 string 类型")
	}
	e.Action = r.Str

	r = gr.Get("params")
	for _, v := range r.Array() {
		e.Params = append(e.Params, v.String())
	}

	r = gr.Get("return")
	if r.Type != gjson.String {
		return e, errors.New("return 字段必须是 string 类型")
	}
	switch r.Str {
	default:
		return e, errors.New("return 字段只能是 string/number/bool/json 中的一种")
	case "string":
		e.Return = StringType
	case "number":
		e.Return = NumberType
	case "bool":
		e.Return = BoolType
	case "json":
		e.Return = JsonType
	}

	r = gr.Get("target")
	if r.Type == gjson.Null {
		e.Target = ""
	} else if r.Type == gjson.String {
		e.Target = r.Str
	} else {
		return e, errors.New("target 字段必须是 null 或 string 类型")
	}

	r = gr.Get("fields")
	for _, v := range r.Array() {
		f, err := parseField(v)
		if err != nil {
			return e, err
		}
		e.Fields = append(e.Fields, f)
	}

	return e, nil
}

func parseField(gr gjson.Result) (Field, error) {
	var t gjson.Result
	var f Field

	t = gr.Get("key")
	if t.Type != gjson.String {
		return f, errors.New("key 必须是 string 类型")
	}
	f.Key = t.Str

	t = gr.Get("type")
	switch t.Str {
	default:
		return f, errors.New("type 定义不正确，必须是 string/number/bool/json 中的一种")
	case "string":
		f.Type = StringType
	case "number":
		f.Type = NumberType
	case "bool":
		f.Type = BoolType
	case "json":
		f.Type = JsonType
	}

	t = gr.Get("target")
	if t.Type == gjson.String {
		f.Target = t.Str
		// if `target` exists, means this is final step, return
		return f, nil
	}

	t = gr.Get("extract")
	if t.Type == gjson.JSON {
		e, err := parseExtract(t)
		if err != nil {
			return f, err
		}
		f.Extract = e
		// if `extract` exists, return
		return f, nil
	}

	t = gr.Get("fields")
	if t.IsArray() {
		for _, v := range t.Array() {
			child, err := parseField(v)
			if err != nil {
				return f, err
			}
			f.Fields = append(f.Fields, child)
		}
		return f, nil
	}

	// check config error
	return f, errors.New("target/extract/fields 这三个字段必须存在一个")
}
