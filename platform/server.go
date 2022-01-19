package platform

import (
	"context"
	"fmt"
	"github.com/meroxa/valve"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/meroxa/funtime/proto"
	"github.com/oklog/run"
	"google.golang.org/grpc"
)

type ProtoWrapper struct {
	ProcessMethod func(context.Context, *proto.ProcessRecordRequest) (*proto.ProcessRecordResponse, error)
}

func (pw ProtoWrapper) Process(ctx context.Context, record *proto.ProcessRecordRequest) (*proto.ProcessRecordResponse, error) {
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

func wrapFrameworkFunc(f func([]valve.Record) ([]valve.Record, []valve.RecordWithError)) func(ctx context.Context, record *proto.ProcessRecordRequest) (*proto.ProcessRecordResponse, error) {
	return func(ctx context.Context, req *proto.ProcessRecordRequest) (*proto.ProcessRecordResponse, error) {
		rr, rre := f(protoRecordToValveRecord(req))
		if rre != nil {
			// TODO: handle
		}
		return valveRecordToProto(rr), nil
	}
}

func protoRecordToValveRecord(req *proto.ProcessRecordRequest) []valve.Record {
	var rr []valve.Record

	for _, pr := range req.Records {
		vr := valve.Record{
			Key:       pr.GetKey(),
			Payload:   valve.Payload(pr.GetValue()),
			Timestamp: time.Unix(pr.GetTimestamp(), 0),
		}
		rr = append(rr, vr)
	}

	return rr
}

func valveRecordToProto(records []valve.Record) *proto.ProcessRecordResponse {
	var prr []*proto.Record
	for _, vr := range records {
		pr := proto.Record{
			Key:       vr.Key,
			Value:     string(vr.Payload),
			Timestamp: vr.Timestamp.Unix(),
		}
		prr = append(prr, &pr)
	}
	return &proto.ProcessRecordResponse{Records: prr}
}
