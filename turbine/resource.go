package turbine

import "github.com/meroxa/turbine-go/turbine/core"

type Resource interface {
	Records(collection string, cfg ConnectionOptions) (Records, error)
	Write(records Records, collection string) error
	WriteWithConfig(records Records, collection string, cfg ConnectionOptions) error
}

// The following aliases exist for backward compatibility.
// In previous version of turbine-go they were used as
// actual type name.

// Deprecated: Use ConnectionOptions instead
type ResourceConfigs = ConnectionOptions

// Deprecated: Use ConnectionOption instead
type ResourceConfig = ConnectionOption

type ConnectionOption struct {
	Field string
	Value string
}

type ConnectionOptions []ConnectionOption

func (cfg ConnectionOptions) ToProto() *core.Configs {
	conf := []*core.Config{}
	for i := range cfg {
		conf = append(conf,
			&core.Config{
				Field: cfg[i].Field,
				Value: cfg[i].Value,
			})
	}
	return &core.Configs{
		Config: conf,
	}
}

func (cfg ConnectionOptions) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	for _, rc := range cfg {
		m[rc.Field] = rc.Value
	}

	return m
}
