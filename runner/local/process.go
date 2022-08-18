package local

import (
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/meroxa/turbine-go/models"
	"reflect"
	"strings"
	"unsafe"
)

func (r Runner) Process(rr models.Records, fn models.BeamFunction) models.Records {
	// use reflection to access intentionally hidden fields
	inVal := reflect.ValueOf(&rr).Elem().FieldByName("pcol")

	// hack to create reference that can be accessed
	in := reflect.NewAt(inVal.Type(), unsafe.Pointer(inVal.UnsafeAddr())).Elem()
	inRR := in.Interface().(beam.PCollection)

	funcName := strings.ToLower(reflect.TypeOf(fn).Name())
	r.scope = r.scope.Scope(funcName)

	out := beam.ParDo(r.scope, fn, inRR)

	return models.NewRecordsWithPcol(out)
}
