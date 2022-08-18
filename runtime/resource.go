package runtime

type Resource struct {
	Source      bool
	Destination bool
	Name        string
	Collection  string
}
type ResourceConfig struct {
	Field string
	Value string
}

type ResourceConfigs []ResourceConfig

func (cfg ResourceConfigs) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	for _, rc := range cfg {
		m[rc.Field] = rc.Value
	}

	return m
}
