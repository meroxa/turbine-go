package server

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/meroxa/turbine-core/v2/proto/process/v2"
	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"

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

func (s *server) Source(n, p string, o ...sdk.Option) (sdk.Source, error) {
	return &source{}, nil
}

func (s *server) SourceWithContext(ctx context.Context, n, p string, o ...sdk.Option) (sdk.Source, error) {
	return &source{}, nil
}

func (s *server) Destination(n, p string, o ...sdk.Option) (sdk.Destination, error) {
	return &destination{}, nil
}

func (s *server) DestinationWithContext(ctx context.Context, n, p string, option ...sdk.Option) (sdk.Destination, error) {
	return &destination{}, nil
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

	processv2.RegisterProcessorServiceServer(s.g, &function{process: fn.Process})
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

func (s *server) Process(rs sdk.Records, fn sdk.Function) (sdk.Records, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var fnName string

	switch reflect.ValueOf(fn).Kind() {
	case reflect.Ptr:
		fnName = strings.ToLower(reflect.TypeOf(fn).Elem().Name())
	default:
		fnName = strings.ToLower(reflect.TypeOf(fn).Name())
	}

	s.fns[fnName] = fn

	return sdk.Records{}, nil
}

func (s *server) ProcessWithContext(_ context.Context, rs sdk.Records, fn sdk.Function) (sdk.Records, error) {
	return s.Process(rs, fn)
}

func funcNames(fns map[string]sdk.Function) string {
	var names []string
	for k := range fns {
		names = append(names, k)
	}

	sort.Strings(names)
	return strings.Join(names, ", ")
}
