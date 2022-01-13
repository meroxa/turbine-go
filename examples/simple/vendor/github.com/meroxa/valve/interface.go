package valve

type App interface {
	Run(Valve) error
}

type Valve interface {
	Resources(string) (Resource, error)
	Process(Records, Function) (Records, RecordsWithErrors)
}