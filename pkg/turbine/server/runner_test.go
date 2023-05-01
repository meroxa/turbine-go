package server

import (
	"context"
	"fmt"
	"math/rand"
	"syscall"
	"testing"
	"time"
	"errors"

	sdk "github.com/meroxa/turbine-go/pkg/turbine"
)

type fakeApp struct{
	wantErr error
}

func (f *fakeApp) Run(t sdk.Turbine) error {
	t.Process(sdk.Records{}, processor{})
	return f.wantErr
}

func Test_Run(t *testing.T) {
	tests := []struct {
		name string
		addr string
		sig  syscall.Signal
		app sdk.App
		wantErr error
	}{
		{
			name: "successfully started",
			addr: fmt.Sprintf(":%d", 51000+rand.Intn(100)),
			app: &fakeApp{},
			sig:  syscall.SIGTERM,
		},
		{
			name: "started then interrupted",
			addr: fmt.Sprintf(":%d", 51000+rand.Intn(100)),
			app: &fakeApp{},
			sig:  syscall.SIGINT,
		},
		{
			name: "returns an error",
			app: &fakeApp{wantErr: errors.New("failure")},
			wantErr: errors.New("failure"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantErr != nil {
				err := Run(context.Background(), tc.app, tc.addr, "processor")
				if err == nil {
					t.Fatalf("expected error %s", tc.wantErr)
				}
				if err.Error() != tc.wantErr.Error() {
					t.Fatalf("want %s, got %s", tc.wantErr.Error(), err.Error())
				}
				return
			}

			exitchan := make(chan bool)
			ready := make(chan bool)

			go func() {
				err := Run(context.Background(), tc.app, tc.addr, "processor")
				if err != nil {
					panic(err)
				}
				exitchan <- true
			}()

			go waitForService(tc.addr, ready)
			<-ready

			syscall.Kill(syscall.Getpid(), tc.sig)

			select {
			case <-exitchan:
				break
			case <-time.After(2 * time.Second):
				t.Fatalf("failed to stop server")
				break
			}
		})
	}
}
