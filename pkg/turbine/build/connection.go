package build

import (
	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	sdk "github.com/meroxa/turbine-go/v2/pkg/turbine"
)

func connectionOptions(opts sdk.ConnectionOptions) *pb.Configs {
	conf := []*pb.Config{}
	for _, co := range opts {
		conf = append(conf,
			&pb.Config{
				Field: co.Field,
				Value: co.Value,
			})
	}
	return &pb.Configs{
		Config: conf,
	}
}
