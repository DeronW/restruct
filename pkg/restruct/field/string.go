package field

import (
	"errors"
	"mss/restruct/pkg/restruct/rule"
	"strconv"

	"github.com/tidwall/gjson"
)

type String struct {
	v     string
	valid bool
}

func NewString(s string) String {
	return String{v: s, valid: true}
}

func (s *String) Value() string {
	return s.v
}

func (s *String) Valid() bool {
	return s.valid
}

func (s *String) jsonParse(fs []rule.Field) ([]Node, error) {
	if !gjson.Valid(s.Value()) {
		return []Node{}, errors.New("json 格式不正确")
	}
	v := gjson.Parse(s.Value())
	return extractJson(v, fs)
}

func (s *String) xmlParse() []Node {
	return nil
}

func (s *String) kvParse(exp string) []Node {
	return nil
}

func (s *String) regexpParse(exp string) []Node {
	return nil
}

func (s *String) mask(exp string) String {
	return NewString("")
}

func (s *String) replace(exp string, replacement string) String {
	if !s.valid {
		return String{v: "", valid: false}
	}
	return NewString("")
}

func (s *String) lstrip(sub string, pos int) String {
	return NewString("")
}

func (s *String) rstrip(exp string, pos int) String {
	return NewString("")
}

func (s *String) toBool() Bool {
	if !s.valid {
		return Bool{v: false, valid: false}
	}

	var v bool = true
	if s.v == "" {
		v = false
	}
	return Bool{v: v, valid: true}
}

func (s *String) toNumber() Number {
	if !s.valid {
		return Number{v: 0, valid: false}
	}

	v, err := strconv.ParseFloat(s.v, 64)

	if err != nil {
		return Number{v: 0, valid: false}
	}

	return Number{v: v, valid: true}
}

func extractJson(r gjson.Result, fs []rule.Field) (ns []Node, err error) {
	for _, f := range fs {
		t := r.Get(f.Key)
		n := Node{Key: f.Key, Type: rule.NullType, Target: f.Target}

		if t.Type == gjson.String && f.Type == rule.StringType {
			n.Type = rule.StringType
			n.Str = NewString(t.Str)
		} else if t.Type == gjson.Number && f.Type == rule.NumberType {
			n.Type = rule.NumberType
			n.Num = NewNumber(t.Num)
		} else if t.IsBool() && f.Type == rule.BoolType {
			n.Type = rule.BoolType
			n.Bool = NewBool(t.Bool())
		} else if t.Type == gjson.JSON && f.Type == rule.JsonType {
			n.Type = rule.JsonType
			n.Children, err = extractJson(t, f.Fields)
		}
		ns = append(ns, n)
	}

	return ns, err
}
