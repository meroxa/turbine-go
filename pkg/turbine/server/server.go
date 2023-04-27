package server

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"strings"
	"sync"

	sdk "github.com/meroxa/turbine-go/pkg/turbine"
	pb "github.com/meroxa/turbine-go/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var _ sdk.Turbine = (*server)(nil)

type server struct {
	mu  sync.Mutex
	g   *grpc.Server
	fns map[string]sdk.Function
}

func NewServer() *server {
	return &server{
		fns: make(map[string]sdk.Function),
		g:   grpc.NewServer(),
	}
}

func (s *server) Listen(addr, name string) error {
	fn, ok := s.fns[name]
	if !ok {
		return fmt.Errorf("cannot find function %q, available functions: %s", name, funcNames(s.fns))
	}

	pb.RegisterFunctionServer(s.g, &function{process: fn.Process})
	healthpb.RegisterHealthServer(s.g, func() healthpb.HealthServer {
		h := health.NewServer()
		h.SetServingStatus("function", healthpb.HealthCheckResponse_SERVING)
		return h
	}())

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return s.g.Serve(listener)
}

func (s *server) GracefulStop() {
	s.g.GracefulStop()
}

func (s *server) Resources(n string) (sdk.Resource, error) {
	return &resource{}, nil
}

func (s *server) ResourcesWithContext(ctx context.Context, n string) (sdk.Resource, error) {
	return &resource{}, nil
}

func (s *server) Process(rs sdk.Records, fn sdk.Function) (sdk.Records, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fnName := strings.ToLower(reflect.TypeOf(fn).Name())
	s.fns[fnName] = fn

	return sdk.Records{}, nil
}

func (s *server) ProcessWithContext(ctx context.Context, rs sdk.Records, fn sdk.Function) (sdk.Records, error) {
	return sdk.Records{}, nil
}

func funcNames(fns map[string]sdk.Function) string {
	var names []string
	for k := range fns {
		names = append(names, k)
	}

	return strings.Join(names, ", ")
}
