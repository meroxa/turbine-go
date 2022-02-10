package turbine

type App interface {
	Run(Valve) error
}

type Valve interface {
	Resources(string) (Resource, error)
	Process(Records, Function) (Records, RecordsWithErrors)
	RegisterSecret(string) error
}
