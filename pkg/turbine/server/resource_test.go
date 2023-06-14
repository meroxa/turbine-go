package server

import (
	"context"
	"testing"

	sdk "github.com/meroxa/turbine-go/v2/pkg/turbine"
)

func Test_Records(t *testing.T) {
	r := resource{}
	rs, err := r.Records("collection", sdk.ConnectionOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(rs.Records); got != 0 {
		t.Fatalf("wanted 0 records, got %d", got)
	}
}

func Test_Write(t *testing.T) {
	r := resource{}

	if err := r.Write(sdk.Records{}, "collection"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := r.WriteWithContext(context.Background(), sdk.Records{}, "collection"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := r.WriteWithConfig(sdk.Records{}, "collection", sdk.ConnectionOptions{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
