package turbine

type Resource interface {
	Records(collection string, cfg ResourceConfigs) (Records, error)
	Write(records Records, collection string) error
	WriteWithConfig(records Records, collection string, cfg ResourceConfigs) error
	GetSpecConnectors() []SpecConnector
}

type SpecConnector struct {
	Type       string                 `json:"type"`
	Resource   string                 `json:"resource"`
	Collection string                 `json:"collection"`
	Config     map[string]interface{} `json:"config,omitempty"`
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
