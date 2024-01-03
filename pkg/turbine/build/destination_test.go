package build

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/conduitio/conduit-commons/opencdc"
	"github.com/conduitio/conduit-commons/proto/opencdc/v1"
	"github.com/meroxa/turbine-core/v2/pkg/client/mock"
	"github.com/meroxa/turbine-core/v2/proto/turbine/v2"
	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"
)

func TestBuilder_AddDestination(t *testing.T) {
	tests := []struct {
		desc    string
		setup   func(*testing.T) *builder
		wantErr error
	}{
		{
			desc: "fails to add destination",
			setup: func(t *testing.T) *builder {
				m := mock.NewMockClient(gomock.NewController(t))
				m.EXPECT().AddDestination(gomock.Any(), &turbinev2.AddDestinationRequest{
					Name: "test",
					Plugin: &turbinev2.Plugin{
						Name: "test-plugin",
						Config: map[string]string{
							"key1": "val1",
							"key2": "val2",
						},
					},
				}).Return(nil, errors.New("boom"))
				return &builder{c: m}
			},
			wantErr: errors.New("boom"),
		},
		{
			desc: "adds destination",
			setup: func(t *testing.T) *builder {
				m := mock.NewMockClient(gomock.NewController(t))
				m.EXPECT().AddDestination(gomock.Any(), &turbinev2.AddDestinationRequest{
					Name: "test",
					Plugin: &turbinev2.Plugin{
						Name: "test-plugin",
						Config: map[string]string{
							"key1": "val1",
							"key2": "val2",
						},
					},
				}).Return(&turbinev2.AddDestinationResponse{
					Id: "dest-id",
				}, nil)
				return &builder{c: m}
			},
		},
	}

	for _, tc := range tests {
		require.True(t, t.Run(tc.desc, func(t *testing.T) {
			b := tc.setup(t)
			_, err := b.Destination("test", "test-plugin", sdk.WithPluginConfig(map[string]string{
				"key1": "val1",
				"key2": "val2",
			}))

			if tc.wantErr != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		}))
	}
}

func TestBuilder_Write(t *testing.T) {
	tests := []struct {
		desc    string
		setup   func(*testing.T) *destination
		records sdk.Records
		wantErr error
	}{
		{
			desc: "fails to write records",
			setup: func(t *testing.T) *destination {
				m := mock.NewMockClient(gomock.NewController(t))
				m.EXPECT().WriteRecords(gomock.Any(), &turbinev2.WriteRecordsRequest{
					DestinationID: "dest-id",
					StreamRecords: &turbinev2.StreamRecords{
						StreamName: "from-source",
						Records: []*opencdcv1.Record{
							{
								Operation: opencdcv1.Operation_OPERATION_CREATE,
								Payload:   new(opencdcv1.Change),
							},
						},
					},
				}).Return(nil, errors.New("boom"))
				return &destination{
					id: "dest-id",
					c:  m,
				}
			},
			records: sdk.Records{
				Stream: "from-source",
				Records: []opencdc.Record{
					{
						Operation: opencdc.OperationCreate,
					},
				},
			},
			wantErr: errors.New("boom"),
		},
		{
			desc: "write records",
			setup: func(t *testing.T) *destination {
				m := mock.NewMockClient(gomock.NewController(t))
				m.EXPECT().WriteRecords(gomock.Any(), &turbinev2.WriteRecordsRequest{
					DestinationID: "dest-id",
					StreamRecords: &turbinev2.StreamRecords{
						StreamName: "from-source",
						Records: []*opencdcv1.Record{
							{
								Operation: opencdcv1.Operation_OPERATION_CREATE,
								Payload:   new(opencdcv1.Change),
							},
						},
					},
				}).Return(new(emptypb.Empty), nil) // new(emptypb.Empty)
				return &destination{
					id: "dest-id",
					c:  m,
				}
			},
			records: sdk.Records{
				Stream: "from-source",
				Records: []opencdc.Record{
					{
						Operation: opencdc.OperationCreate,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		require.True(t, t.Run(tc.desc, func(t *testing.T) {
			err := tc.setup(t).Write(tc.records)

			if tc.wantErr != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		}))
	}
}
