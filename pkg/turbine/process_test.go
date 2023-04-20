package turbine

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/turbine-go/pkg/proto/core"
	"github.com/meroxa/turbine-go/pkg/turbine/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ClearAllFunc struct{}

func (t ClearAllFunc) Process(r []Record) []Record {
	return []Record{}
}

func TestProcess(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		turbineMock = mock.NewMockTurbineCore(ctrl)
		now         = timestamppb.Now()
	)

	turbineMock.EXPECT().
		AddProcessToCollection(gomock.Any(), &core.ProcessCollectionRequest{
			Collection: &core.Collection{
				Name:   "name",
				Stream: "stream",
				Records: []*core.Record{
					{
						Key:       "key",
						Value:     []byte("payload"),
						Timestamp: now,
					},
				},
			},
			Process: &core.ProcessCollectionRequest_Process{
				Name: "clearallfunc",
			},
		}).Times(1).
		Return(&core.Collection{
			Name:   "name",
			Stream: "stream",
			Records: []*core.Record{
				{
					Key:       "key",
					Value:     []byte("payload"),
					Timestamp: now,
				},
			},
		}, nil)
	tb := turbine{TurbineCore: turbineMock}

	rs, err := tb.Process(
		Records{
			Name:   "name",
			Stream: "stream",
			Records: []Record{
				{
					Key:       "key",
					Payload:   []byte("payload"),
					Timestamp: now.AsTime(),
				},
			},
		},
		ClearAllFunc{},
	)
	require.NoError(t, err)
	require.Equal(t, Records{Stream: "stream", Records: []Record{}, Name: "name"}, rs)
}
