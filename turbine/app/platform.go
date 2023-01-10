package app

import (
	"context"
	"os"
	"time"

	pb "github.com/meroxa/turbine-go/turbine/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Turbine struct {
	*grpc.ClientConn
	pb.TurbineServiceClient
	recording bool
}

func (t *Turbine) Recording() bool {
	return t.recording
}

func New(ctx context.Context, recording bool) (*Turbine, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(
		ctx,
		os.Getenv("TURBINE_CORE_SERVER"),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Turbine{
		ClientConn:           conn,
		TurbineServiceClient: pb.NewTurbineServiceClient(conn),
		recording:            recording,
	}, nil
}
