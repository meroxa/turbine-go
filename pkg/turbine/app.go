package turbine

type App interface {
	Run(Turbine) error
}
