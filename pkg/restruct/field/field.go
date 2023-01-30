package field

import (
	"errors"
	"mss/restruct/pkg/restruct/rule"
	"strconv"
)

type Node struct {
	Key      string
	Target   string
	Type     rule.FieldType
	Str      String
	Num      Number
	Bool     Bool
	Children []Node
}

func StringExtract(s String, e rule.Extract) ([]Node, error) {
	if !s.Valid() {
		return []Node{}, errors.New("字段解析失败")
	}
	var ns []Node
	var err error

	switch e.Action {
	default:
		err = errors.New("提取方法不存在")
	case "json":
		ns, err = s.jsonParse(e.Fields)
	case "xml":
		ns = s.xmlParse()
	case "key-value":
		ns = s.kvParse(e.Params[0])
	case "regexp":
		ns = s.regexpParse(e.Params[0])
	case "lstrip":
		exp := e.Params[0]
		pos, err := strconv.ParseInt(e.Params[1], 10, 32)
		if err == nil {
			ns = []Node{{
				Type: rule.StringType,
				Str:  s.lstrip(exp, int(pos)),
			}}
		}
	case "rstrip":
		exp := e.Params[0]
		pos, err := strconv.ParseInt(e.Params[1], 10, 32)
		if err == nil {
			ns = []Node{{
				Type: rule.StringType,
				Str:  s.rstrip(exp, int(pos)),
			}}
		}
	case "mask":
		ns = []Node{{
			Type: rule.StringType,
			Str:  s.mask(e.Params[0]),
		}}
	case "replace":
		ns = []Node{{
			Type: rule.StringType,
			Str:  s.replace(e.Params[0], e.Params[1]),
		}}
	case "toNumber":
		ns = []Node{{Type: rule.NumberType, Num: s.toNumber()}}
	case "toBool":
		ns = []Node{{Type: rule.BoolType, Bool: s.toBool()}}
	}
	return ns, err
}

func NumberExtract(n Number, e rule.Extract) (ns []Node, err error) {
	if n.Valid() {
		return
	}
	switch e.Action {
	default:
		err = errors.New("提取方法不存在")
	case "toString":
		ns = append(ns, Node{Type: rule.StringType, Str: n.toString()})
	case "toBool":
		ns = append(ns, Node{Type: rule.BoolType, Bool: n.toBool()})
	}
	return
}

func BoolExtract(b Bool, e rule.Extract) ([]Node, error) {
	switch e.Action {
	default:
		return nil, errors.New("提取方法不存在")
	case "toString":
		return []Node{{Type: rule.StringType, Str: b.toString()}}, nil
	case "toNumber":
		return []Node{{Type: rule.NumberType, Num: b.toNumber()}}, nil
	}
}
