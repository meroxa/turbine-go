package app

import (
	"time"

	"github.com/meroxa/turbine-go/pkg/proto/core"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Records struct {
	Stream  string
	Records []Record
	Name    string
}

func NewRecords(c *core.Collection) Records {
	rs := []Record{}
	for _, r := range c.Records {
		rs = append(rs,
			Record{
				Key:       r.Key,
				Payload:   r.Value,
				Timestamp: r.Timestamp.AsTime(),
			},
		)
	}

	return Records{
		Stream:  c.Stream,
		Records: rs,
		Name:    c.Name,
	}
}

func (rs *Records) ToProto() *core.Collection {
	rds := []*core.Record{}
	for _, r := range rs.Records {
		rds = append(rds,
			&core.Record{
				Key:       r.Key,
				Value:     r.Payload,
				Timestamp: timestamppb.New(r.Timestamp),
			})
	}
	return &core.Collection{
		Stream:  rs.Stream,
		Records: rds,
		Name:    rs.Name,
	}
}

type Record struct {
	Key       string
	Payload   Payload
	Timestamp time.Time
}

type Payload []byte
