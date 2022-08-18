package runtime

type Processor struct {
	Name     string
	Inputs   map[string]string
	Pipeline Pipeline
}
