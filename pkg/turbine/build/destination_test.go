package build

/*

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/meroxa/turbine-core/v2/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/client/mock"
	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"
)
func TestDestination(t *testing.T) {
	var (
		ctrl       = gomock.NewController(t)
		clientMock = mock.NewMockClient(ctrl)

		resourceName = "pg"
		coreSource   = &pb.Destination{
			Name: resourceName,
		}
	)

	clientMock.EXPECT().
		GetDestination(gomock.Any(), &pb.GetDestinationRequest{
			Name: resourceName,
		}).Times(1).
		Return(coreSource, nil)

	b := builder{Client: clientMock}
	r, err := b.Destination(resourceName)
	require.NoError(t, err)
	require.Equal(t, r, &destination{
		Destination: coreSource,
		Client:      clientMock,
	})
}

func TestRecords(t *testing.T) {
	var (
		ctrl       = gomock.NewController(t)
		clientMock = mock.NewMockClient(ctrl)

		r = &source{
			Source: &pb.Source{Name: "pg"},
			Client: clientMock,
		}
		collection = "accounts"
	)

	clientMock.EXPECT().
		ReadCollection(gomock.Any(), &pb.ReadCollectionRequest{
			Source:     r.Source,
			Collection: collection,
			Configs: &pb.Configs{
				Config: []*pb.Config{
					{
						Field: "conf",
						Value: "conf_val",
					},
				},
			},
		}).Times(1).
		Return(&pb.Collection{
			Name:   "name",
			Stream: "stream",
		}, nil)

	c, err := r.Read()

	require.NoError(t, err)
	require.Equal(t, c, sdk.Records{
		Name:    "name",
		Stream:  "stream",
		Records: []sdk.Record{},
	})
}

func TestWrite(t *testing.T) {
	var (
		ctrl       = gomock.NewController(t)
		clientMock = mock.NewMockClient(ctrl)

		r = &destination{
			Destination: &pb.Destination{Name: "pg"},
			Client:      clientMock,
		}
		targetCollection = "copy"
	)

	clientMock.EXPECT().
		WriteCollectionToDestination(gomock.Any(), &pb.WriteCollectionRequest{
			Destination: r.Destination,
			SourceCollection: &pb.Collection{
				Name:    "name",
				Stream:  "stream",
				Records: []*pb.Record{},
			},
			DestinationCollection: targetCollection,
			Configs: &pb.Configs{
				Config: []*pb.Config{},
			},
		}).Times(1).
		Return(&emptypb.Empty{}, nil)

	require.NoError(t, r.Write(sdk.Records{
		Name:    "name",
		Stream:  "stream",
		Records: []sdk.Record{},
	}))
}

*/
