package build

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	//"github.com/conduitio/conduit-commons/proto/opencdc/v1"
	"github.com/meroxa/turbine-core/v2/pkg/client/mock"
	"github.com/meroxa/turbine-core/v2/proto/turbine/v2"
	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"
)

func TestBuilder_AddSource(t *testing.T) {
	tests := []struct {
		desc    string
		setup   func(*testing.T) *builder
		wantErr error
	}{
		{
			desc: "fails to add destination",
			setup: func(t *testing.T) *builder {
				m := mock.NewMockClient(gomock.NewController(t))
				m.EXPECT().AddSource(gomock.Any(), &turbinev2.AddSourceRequest{
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
				m.EXPECT().AddSource(gomock.Any(), &turbinev2.AddSourceRequest{
					Name: "test",
					Plugin: &turbinev2.Plugin{
						Name: "test-plugin",
						Config: map[string]string{
							"key1": "val1",
							"key2": "val2",
						},
					},
				}).Return(&turbinev2.AddSourceResponse{
					Id: "dest-id",
				}, nil)
				return &builder{c: m}
			},
		},
	}

	for _, tc := range tests {
		require.True(t, t.Run(tc.desc, func(t *testing.T) {
			b := tc.setup(t)
			_, err := b.Source("test", "test-plugin", sdk.WithPluginConfig(map[string]string{
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

func TestBuilder_Read(t *testing.T) {
	records := testRecords(t, "source-stream")

	setup := func(t *testing.T, withErr error) *source {
		m := mock.NewMockClient(gomock.NewController(t))
		call := m.EXPECT().ReadRecords(gomock.Any(), &turbinev2.ReadRecordsRequest{
			SourceStream: "source-stream",
		})
		if withErr != nil {
			call.Return(nil, withErr)
		} else {
			call.Return(&turbinev2.ReadRecordsResponse{
				StreamRecords: &turbinev2.StreamRecords{
					StreamName: "source-stream",
					Records:    protoRecords(t, records),
				},
			}, nil)
		}

		return &source{
			streamName: "source-stream",
			c:          m,
		}
	}

	tests := []struct {
		desc    string
		source  *source
		wantErr error
	}{
		{
			desc:    "fails to read records",
			source:  setup(t, errors.New("boom")),
			wantErr: errors.New("boom"),
		},
		{
			desc:   "read records",
			source: setup(t, nil),
		},
	}

	for _, tc := range tests {
		require.True(t, t.Run(tc.desc, func(t *testing.T) {
			rs, err := tc.source.Read()

			if tc.wantErr != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.wantErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, rs, records)
			}
		}))
	}
}
