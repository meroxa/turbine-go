package server

import (
	"context"

	"github.com/conduitio/conduit-commons/opencdc"
	"github.com/conduitio/conduit-commons/proto/opencdc/v1"
	"github.com/meroxa/turbine-core/v2/proto/process/v2"
)

type function struct {
	processv2.UnimplementedProcessorServiceServer

	process func([]opencdc.Record) []opencdc.Record
}

func (f *function) Process(ctx context.Context, req *processv2.ProcessRequest) (*processv2.ProcessResponse, error) {
	// unmarshal proto to opencdc records
	in := make([]opencdc.Record, len(req.Records))

	for i, pr := range req.Records {
		if err := in[i].FromProto(pr); err != nil {
			return nil, err
		}
	}

	// pass records to function for processing
	processed := f.process(in)

	// marsha opencdc records to proto
	out := make([]*opencdcv1.Record, len(processed))

	for i, r := range processed {
		out[i] = &opencdcv1.Record{}
		if err := r.ToProto(out[i]); err != nil {
			return nil, err
		}
	}

	return &processv2.ProcessResponse{Records: out}, nil
}
