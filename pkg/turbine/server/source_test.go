package server

import (
	"context"
	"testing"

	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"
)

func Test_Records(t *testing.T) {
	r := source{}
	rs, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(rs.Records); got != 0 {
		t.Fatalf("wanted 0 records, got %d", got)
	}
}

func Test_Write(t *testing.T) {
	r := destination{}

	if err := r.Write(sdk.Records{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := r.WriteWithContext(context.Background(), sdk.Records{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
