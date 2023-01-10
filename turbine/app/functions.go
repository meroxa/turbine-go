package app

import (
	"context"
	"reflect"
	"strings"

	"github.com/meroxa/turbine-go/turbine"
	"github.com/meroxa/turbine-go/turbine/core"
)

func (t *Turbine) Process(rr turbine.Records, fn turbine.Function) (turbine.Records, error) {
	c, err := t.AddProcessToCollection(context.TODO(), &core.ProcessCollectionRequest{
		Process: &core.ProcessCollectionRequest_Process{
			Name: strings.ToLower(reflect.TypeOf(fn).Name()),
		},
		Collection: rr.ToProto(),
	})
	if err != nil {
		return turbine.Records{}, err
	}

	// // use reflection to access intentionally hidden fields
	// inVal := reflect.ValueOf(&rr).Elem().FieldByName("Records")

	// // hack to create reference that can be accessed
	// in := reflect.NewAt(inVal.Type(), unsafe.Pointer(inVal.UnsafeAddr())).Elem()
	// inRR := in.Interface().([]turbine.Record)

	// rawOut := fn.Process(inRR)
	// rr.Records = rawOut
	records := turbine.NewRecords(c)
	if t.recording {
		return records, nil
	}

	// // use reflection to access intentionally hidden fields
	// inVal := reflect.ValueOf(&rr).Elem().FieldByName("Records")

	// // hack to create reference that can be accessed
	// in := reflect.NewAt(inVal.Type(), unsafe.Pointer(inVal.UnsafeAddr())).Elem()
	// inRR := in.Interface().([]turbine.Record)

	rawOut := fn.Process(records.Records)
	records.Records = rawOut
	return records, nil

}

// func (t *Turbine) process(fn turbine.Function, c *core.Collection) turbine.Records {
// 	records := turbine.NewRecords(c)
// 	if t.recording {
// 		return records
// 	}

// 	// // use reflection to access intentionally hidden fields
// 	// inVal := reflect.ValueOf(&rr).Elem().FieldByName("Records")

// 	// // hack to create reference that can be accessed
// 	// in := reflect.NewAt(inVal.Type(), unsafe.Pointer(inVal.UnsafeAddr())).Elem()
// 	// inRR := in.Interface().([]turbine.Record)

// 	rawOut := fn.Process(records.Records)
// 	records.Records = rawOut
// 	return records

// }
