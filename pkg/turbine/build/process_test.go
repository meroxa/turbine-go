package build

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/client/mock"

	sdk "github.com/meroxa/turbine-go/v2/pkg/turbine"
)

type ClearAllFunc struct{}

func (t ClearAllFunc) Process(r []sdk.Record) []sdk.Record {
	return []sdk.Record{}
}

func TestProcess(t *testing.T) {
	var (
		ctrl       = gomock.NewController(t)
		clientMock = mock.NewMockClient(ctrl)
		now        = timestamppb.Now()
	)

	clientMock.EXPECT().
		AddProcessToCollection(gomock.Any(), &pb.ProcessCollectionRequest{
			Collection: &pb.Collection{
				Name:   "name",
				Stream: "stream",
				Records: []*pb.Record{
					{
						Key:       "key",
						Value:     []byte("payload"),
						Timestamp: now,
					},
				},
			},
			Process: &pb.ProcessCollectionRequest_Process{
				Name: "clearallfunc",
			},
		}).Times(1).
		Return(&pb.Collection{
			Name:   "name",
			Stream: "stream",
			Records: []*pb.Record{
				{
					Key:       "key",
					Value:     []byte("payload"),
					Timestamp: now,
				},
			},
		}, nil)
	b := builder{Client: clientMock}

	rs, err := b.Process(
		sdk.Records{
			Name:   "name",
			Stream: "stream",
			Records: []sdk.Record{
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
	require.Equal(t, sdk.Records{Stream: "stream", Records: []sdk.Record{}, Name: "name"}, rs)
}
