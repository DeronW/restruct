package field

type Bool struct {
	v     bool
	valid bool
}

func NewBool(b bool) Bool {
	return Bool{v: b, valid: true}
}

func (b *Bool) toString() String {
	if !b.valid {
		return String{v: "", valid: false}
	}

	var v string = "false"
	if b.v {
		v = "true"
	}
	return String{v: v, valid: true}
}

func (b *Bool) toNumber() Number {
	if !b.valid {
		return Number{v: 0, valid: false}
	}

	var v float64 = 0
	if b.v {
		v = 1
	}
	return Number{v: v, valid: true}
}

func (b *Bool) Value() bool {
	return b.v
}

func (b *Bool) Valid() bool {
	return b.valid
}
