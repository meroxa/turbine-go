package turbine

import	"github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"

type ConnectionOption struct {
	Field string
	Value string
}

type ConnectionOptions []ConnectionOption

func (cos ConnectionOptions) ToProto() *core.Configs {
	conf := []*core.Config{}
	for _, co := range cos {
		conf = append(conf,
			&core.Config{
				Field: co.Field,
				Value: co.Value,
			})
	}
	return &core.Configs{
		Config: conf,
	}
}

func (cos ConnectionOptions) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	for _, rc := range cos {
		m[rc.Field] = rc.Value
	}

	return m
}
