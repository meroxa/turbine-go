package models

import (
	"context"
)

type Function interface {
	Process(r Record) Record
}

type BeamFunction interface {
	ProcessElement(context.Context, string, func(string))
}
