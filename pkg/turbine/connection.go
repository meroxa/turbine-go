package turbine

type ConnectionOption struct {
	Field string
	Value string
}

type ConnectionOptions []ConnectionOption

func (cos ConnectionOptions) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	for _, rc := range cos {
		m[rc.Field] = rc.Value
	}

	return m
}
