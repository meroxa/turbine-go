package turbine

type App interface {
	Run(Server) error
}
