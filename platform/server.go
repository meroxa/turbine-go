package platform

import (
	"context"
	"fmt"
	"github.com/meroxa/valve"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/meroxa/valve/platform/proto"
	"github.com/oklog/run"
	"google.golang.org/grpc"
)

type ProtoWrapper struct {
	ProcessMethod func(ctx context.Context, record *proto.Record) (*proto.Record, error)
}

func (pw ProtoWrapper) Process(ctx context.Context, record *proto.Record) (*proto.Record, error) {
	return pw.ProcessMethod(ctx, record)
}

func ServeFunc(f valve.Function) error {

	convertedFunc := wrapFrameworkFunc(f.Process)

	fn := struct{ ProtoWrapper }{}
	fn.ProcessMethod = convertedFunc

	addr := os.Getenv("MEROXA_FUNCTION_ADDR")
	if addr == "" {
		return fmt.Errorf("Missing MEROXA_FUNCTION_ADDR env var")
	}

	var g run.Group
	g.Add(run.SignalHandler(context.Background(), syscall.SIGTERM))
	{
		gsrv := grpc.NewServer()
		proto.RegisterFunctionServer(gsrv, fn)

		g.Add(func() error {
			ln, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}

			return gsrv.Serve(ln)
		}, func(err error) {
			gsrv.GracefulStop()
		})
	}

	return g.Run()
}

func wrapFrameworkFunc(f func([]valve.Record) ([]valve.Record, []valve.RecordWithError)) func(ctx context.Context, record *proto.Record) (*proto.Record, error) {
	return func(ctx context.Context, record *proto.Record) (*proto.Record, error) {
		rr, rre := f([]valve.Record{protoRecordToValveRecord(record)})
		if rre != nil {
			// TODO: handle
		}
		return valveRecordToProto(rr[0]), nil
	}
}

func protoRecordToValveRecord(record *proto.Record) valve.Record {
	return valve.Record{
		Key:       record.Key,
		Payload:   valve.Payload(record.Value),
		Timestamp: time.Unix(record.Timestamp, 0),
	}
}

func valveRecordToProto(record valve.Record) *proto.Record {
	return &proto.Record{
		Key:       record.Key,
		Value:     string(record.Payload),
		Timestamp: record.Timestamp.Unix(),
	}

}
