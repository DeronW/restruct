package restruct

import (
	"encoding/json"
	"errors"
	"mss/restruct/pkg/restruct/field"
	"mss/restruct/pkg/restruct/rule"
)

type ValueType string

const (
	String ValueType = "string"
	Null   ValueType = "null"
	Number ValueType = "number"
	Bool   ValueType = "bool"
)

type Value struct {
	Key  string
	Type ValueType
	Str  string
	Num  float64
	Bool bool
}

/*
This is entry point of re struct data. It follow the extract rule tree,
and extract all data by deepth-first traversal, until extract rule is run out.
there are 2 kinds of node in extract tree:
1. extract: mark data trasnfer rules
2. fields: mark extracted fields
each of them could contain children or not, but only field could contain extract.
*/
func Restruct(s string, e rule.Extract) ([]Value, error) {
	node := field.Node{Str: field.NewString(s), Type: rule.StringType}
	nodes, err := travel(node, e)
	return nodes, err
}

func travel(n field.Node, e rule.Extract) ([]Value, error) {
	var err error
	var nodes []field.Node
	values := []Value{}

	if n.Type == rule.StringType {
		nodes, err = field.StringExtract(n.Str, e)
	}
	if n.Type == rule.NumberType {
		nodes, err = field.NumberExtract(n.Num, e)
	}
	if n.Type == rule.BoolType {
		nodes, err = field.BoolExtract(n.Bool, e)
	}

	// got final value directly
	if e.Target != "" {
		if len(nodes) != 1 {
			err = errors.New("解析结果数量不匹配，应该得到一个结果")
		} else {
			values = append(values, node2value(e.Target, nodes[0]))
		}
		return values, err
	}

	// loop get children value
	for i := 0; i < len(e.Fields); i++ {
		field := e.Fields[i]
		node := nodes[i]

		if field.Target != "" {
			values = append(values, node2value(field.Target, node))
			continue
		}

		if field.Extract.Action != "" {
			var vs []Value
			vs, err = travel(node, field.Extract)
			values = append(values, vs...)
			continue
		}

		// result is nested level json
		for j := 0; j < len(field.Fields) && j < len(node.Children); j++ {
			var vs []Value
			vs, err = travelFields(node.Children[j], field.Fields[j])
			values = append(values, vs...)
		}
	}

	return values, err
}

func travelFields(n field.Node, f rule.Field) (vs []Value, err error) {
	if f.Target != "" {
		return []Value{node2value(f.Target, n)}, nil
	}
	if f.Extract.Action != "" {
		vs, err = travel(n, f.Extract)
		return vs, err
	}

	l := len(n.Children)
	if l > len(f.Fields) {
		l = len(f.Fields)
	}
	for i := 0; i < l; i++ {
		var v []Value
		v, err = travelFields(n.Children[i], f.Fields[i])
		vs = append(vs, v...)
	}
	return vs, err
}

func node2value(target string, n field.Node) Value {
	v := Value{Key: target}

	if n.Type == rule.BoolType {
		v.Type = Bool
		v.Bool = n.Bool.Value()
	} else if n.Type == rule.NumberType {
		v.Type = Number
		v.Num = n.Num.Value()
	} else if n.Type == rule.StringType {
		v.Type = String
		v.Str = n.Str.Value()
	} else if n.Type == rule.NullType {
		v.Type = Null
	}

	return v
}

func ToString(values []Value) (string, error) {
	r := map[string]any{}
	for _, i := range values {
		if i.Type == Null {
			r[i.Key] = nil
		}
		if i.Type == Bool {
			r[i.Key] = i.Bool
		}
		if i.Type == Number {
			r[i.Key] = i.Num
		}
		if i.Type == String {
			r[i.Key] = i.Str
		}
	}
	s, err := json.Marshal(r)
	return string(s), err
}
