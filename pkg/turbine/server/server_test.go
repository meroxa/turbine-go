package server

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	sdk "github.com/meroxa/turbine-go/v2/pkg/turbine"
)

func Test_Resources(t *testing.T) {
	s := &server{}
	if got, err := s.Resources("nothing"); got == nil && err != nil {
		t.Fatalf("got %v, error: %v", got, err)
	}
}

func Test_ResourcesWithContext(t *testing.T) {
	s := &server{}
	got, err := s.ResourcesWithContext(context.Background(), "nothing")
	if got == nil && err != nil {
		t.Fatalf("got %v, error: %v", got, err)
	}
}

func Test_Process(t *testing.T) {
	s := &server{
		fns: make(map[string]sdk.Function),
	}

	tests := []struct {
		name string
		proc sdk.Function
	}{
		{
			name: "with concrete processor",
			proc: processor{},
		},
		{
			name: "with ptr processor",
			proc: &processor{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := s.Process(sdk.Records{}, tc.proc)
			if err != nil {
				t.Fatalf("got %v, error: %v", got, err)
			}
			if fn := s.fns["processor"]; fn != tc.proc {
				t.Fatalf("got %v, wanted %v", fn, tc.proc)
			}
		})
	}
}

func Test_ProcessContext(t *testing.T) {
	s := &server{
		fns: make(map[string]sdk.Function),
	}
	proc := processor{}

	got, err := s.ProcessWithContext(context.Background(), sdk.Records{}, proc)
	if err != nil {
		t.Fatalf("got %v, error: %v", got, err)
	}
	if fn := s.fns["processor"]; fn != proc {
		t.Fatalf("got %v, wanted %v", fn, proc)
	}
}

func Test_Listen(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		setup   func() *server
		wantErr error
	}{
		{
			name:    "error on missing function",
			addr:    ":0",
			setup:   func() *server { return NewServer() },
			wantErr: errors.New(`cannot find function "processor", available functions: `),
		},
		{
			name:    "failed to listen on address",
			wantErr: errors.New("listen tcp: address -1: invalid port"),
			setup: func() *server {
				s := NewServer()
				_, _ = s.Process(sdk.Records{}, processor{})
				return s
			},
			addr: ":-1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := tc.setup()
			err := s.Listen(tc.addr, "processor")
			if tc.wantErr != nil && err != nil {
				if tc.wantErr.Error() != err.Error() {
					t.Fatalf("want %s, got %s", tc.wantErr.Error(), err.Error())
				}
			}
			if tc.wantErr != nil && err == nil {
				t.Fatalf("expected error: %s", tc.wantErr.Error())
			}
		})
	}
}

func Test_GracefulStop(t *testing.T) {
	var (
		s        = NewServer()
		addr     = fmt.Sprintf(":%d", 51000+rand.Intn(100))
		exitchan = make(chan bool)
		ready    = make(chan bool)
	)

	if _, err := s.Process(sdk.Records{}, processor{}); err != nil {
		t.Fatalf("got error: %v", err)
	}

	go func() {
		if err := s.Listen(addr, "processor"); err != nil {
			panic(err)
		}
		exitchan <- true
	}()

	go waitForService(addr, ready)

	<-ready
	s.GracefulStop()

	select {
	case <-exitchan:
		break
	case <-time.After(2 * time.Second):
		t.Fatalf("failed to stop server")
		break
	}
}

func Test_funcNames(t *testing.T) {
	tests := []struct {
		name string
		want string
		fns  map[string]sdk.Function
	}{
		{
			name: "concat with zero function",
			fns:  make(map[string]sdk.Function),
		},
		{
			name: "concat two functions",
			fns: map[string]sdk.Function{
				"proc1": processor{},
				"proc2": processor{},
			},
			want: "proc1, proc2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := funcNames(tc.fns); got != tc.want {
				t.Fatalf("got %s, want %s", got, tc.want)
			}
		})
	}
}

type processor struct{}

func (p processor) Process(rs []sdk.Record) []sdk.Record {
	out := []sdk.Record{}
	for _, r := range rs {
		r.Key = r.Key + "+processed"
		out = append(out, r)
	}
	return out
}

func waitForService(addr string, done chan bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	for {
		select {
		case <-time.After(100 * time.Millisecond):
			c, err := net.Dial("tcp", addr)
			if err == nil {
				c.Close()
				done <- true
				return
			}
		case <-ctx.Done():
			panic(ctx.Err())
		}
	}
}
