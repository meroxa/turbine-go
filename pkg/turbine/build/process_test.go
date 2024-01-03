package build

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/conduitio/conduit-commons/opencdc"
	"github.com/conduitio/conduit-commons/proto/opencdc/v1"
	"github.com/meroxa/turbine-core/v2/pkg/client/mock"
	"github.com/meroxa/turbine-core/v2/proto/turbine/v2"
	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"
)

func TestBuilder_Process(t *testing.T) {
	setup := func(t *testing.T, withProc bool, withErr error) *builder {
		t.Helper()
		sr := &turbinev2.StreamRecords{
			StreamName: "stream1",
			Records:    protoRecords(t, testRecords(t, "stream1")),
		}

		m := mock.NewMockClient(gomock.NewController(t))
		call := m.EXPECT().
			ProcessRecords(gomock.Any(), &turbinev2.ProcessRecordsRequest{
				StreamRecords: sr,
				Process: &turbinev2.ProcessRecordsRequest_Process{
					Name: "testreplacer",
				},
			})

		if withErr != nil {
			call.Return(nil, errors.New("boom"))
		} else {
			call.Return(&turbinev2.ProcessRecordsResponse{
				StreamRecords: sr,
			}, nil)
		}

		return &builder{c: m, runProcess: withProc}
	}

	tests := []struct {
		desc            string
		input, expected sdk.Records
		builder         *builder
		wantErr         error
	}{
		{
			desc:     "running processor",
			input:    testRecords(t, "stream1"),
			expected: processedRecords(t, "stream1"),
			builder:  setup(t, true, nil),
		},
		{
			desc:     "without processor",
			input:    testRecords(t, "stream1"),
			expected: testRecords(t, "stream1"),
			builder:  setup(t, false, nil),
		},
		{
			desc:     "process records error",
			input:    testRecords(t, "stream1"),
			expected: testRecords(t, "stream1"),
			builder:  setup(t, true, errors.New("boom")),
			wantErr:  errors.New("boom"),
		},
	}

	for _, tc := range tests {
		require.True(t, t.Run(tc.desc, func(t *testing.T) {
			rs, err := tc.builder.Process(tc.input, testReplacer{})

			if tc.wantErr != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.wantErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, rs, tc.expected)
			}
		}))
	}
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
