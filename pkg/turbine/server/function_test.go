package server

import (
	"context"
	"testing"

	"github.com/conduitio/conduit-commons/opencdc"
	"github.com/conduitio/conduit-commons/proto/opencdc/v1"
	"github.com/meroxa/turbine-core/v2/proto/process/v2"
	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"
	"github.com/stretchr/testify/require"
)

func Test_function_Process(t *testing.T) {
	fn := function{process: testReplacer{}.Process}

	in := protoRecords(t, testRecords(t, "stream-1"))
	expected := protoRecords(t, processedRecords(t, "stream-1"))

	resp, err := fn.Process(context.Background(), &processv2.ProcessRequest{Records: in})
	require.NoError(t, err)

	require.Equal(t, resp.Records, expected)
}

type testReplacer struct{}

func (r testReplacer) Process(rs []opencdc.Record) []opencdc.Record {
	out := make([]opencdc.Record, len(rs))
	for i, record := range rs {
		out[i] = record.Clone()
		out[i].Payload.Before = out[i].Payload.After.Clone()

		data := out[i].Payload.After.Clone().(opencdc.StructuredData)
		data["replaced"] = "true"

		out[i].Payload.After = data
	}

	return out
}

func protoRecords(t *testing.T, rs sdk.Records) []*opencdcv1.Record {
	t.Helper()

	protoRecords := make([]*opencdcv1.Record, len(rs.Records))
	for i, r := range rs.Records {
		protoRecords[i] = &opencdcv1.Record{}
		require.NoError(t, r.ToProto(protoRecords[i]))
	}

	return protoRecords
}

func processedRecords(t *testing.T, stream string) sdk.Records {
	t.Helper()

	return sdk.Records{
		Stream: stream,
		Records: []opencdc.Record{
			{
				Position:  opencdc.Position("one"),
				Operation: opencdc.OperationCreate,
				Metadata:  opencdc.Metadata{"meta": "data"},
				Key:       opencdc.RawData("magic"),
				Payload: opencdc.Change{
					Before: opencdc.StructuredData{
						"two": "three",
					},
					After: opencdc.StructuredData{
						"two":      "three",
						"replaced": "true",
					},
				},
			},
		},
	}
}

func testRecords(t *testing.T, stream string) sdk.Records {
	t.Helper()
	return sdk.Records{
		Stream: stream,
		Records: []opencdc.Record{
			{
				Position:  opencdc.Position("one"),
				Operation: opencdc.OperationCreate,
				Metadata:  opencdc.Metadata{"meta": "data"},
				Key:       opencdc.RawData("magic"),
				Payload: opencdc.Change{
					Before: opencdc.StructuredData{
						"one": "two",
					},
					After: opencdc.StructuredData{
						"two": "three",
					},
				},
			},
		},
	}
}
