package field

import "fmt"

type Number struct {
	v     float64
	valid bool
}

func NewNumber(n float64) Number {
	return Number{v: n, valid: true}
}

func (n *Number) toString() String {
	if !n.valid {
		return String{v: "", valid: false}
	}
	return String{v: fmt.Sprintf("%f", n.v), valid: true}
}

func (n *Number) toBool() Bool {
	if !n.valid {
		return Bool{v: false, valid: false}
	}

	var b bool = false
	if n.v != 0 {
		b = true
	}
	return Bool{v: b, valid: true}
}

func (n *Number) Value() float64 {
	return n.v
}

func (n *Number) Valid() bool {
	return n.valid
}
