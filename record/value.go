package record

type Value struct {
	value interface{}
}

func (v Value) String() string {
	return v.value.(string)
}

func (v Value) Int() int {
	return v.value.(int)
}

func (v Value) Float() float64 {
	return v.value.(float64)
}
func (v Value) Bytes() []byte {
	return v.value.([]byte)
}
