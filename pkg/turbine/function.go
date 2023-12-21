package turbine

import "github.com/conduitio/conduit-commons/opencdc"

type Function interface {
	Process(r []opencdc.Record) []opencdc.Record
}
