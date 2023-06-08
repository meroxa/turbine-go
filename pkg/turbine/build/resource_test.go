package build

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/client/mock"
	sdk "github.com/meroxa/turbine-go/v2/pkg/turbine"
)

func TestResources(t *testing.T) {
	var (
		ctrl       = gomock.NewController(t)
		clientMock = mock.NewMockClient(ctrl)

		resourceName = "pg"
		coreResource = &pb.Resource{
			Name: resourceName,
		}
	)

	clientMock.EXPECT().
		GetResource(gomock.Any(), &pb.GetResourceRequest{
			Name: resourceName,
		}).Times(1).
		Return(coreResource, nil)

	b := builder{Client: clientMock}
	r, err := b.Resources(resourceName)
	require.NoError(t, err)
	require.Equal(t, r, &resource{
		Resource: coreResource,
		Client:   clientMock,
	})
}

func TestRecords(t *testing.T) {
	var (
		ctrl       = gomock.NewController(t)
		clientMock = mock.NewMockClient(ctrl)

		r = &resource{
			Resource: &pb.Resource{Name: "pg"},
			Client:   clientMock,
		}
		collection = "accounts"
	)

	clientMock.EXPECT().
		ReadCollection(gomock.Any(), &pb.ReadCollectionRequest{
			Resource:   r.Resource,
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

	c, err := r.Records(collection, sdk.ConnectionOptions{{
		Field: "conf",
		Value: "conf_val",
	}})
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

		r = &resource{
			Resource: &pb.Resource{Name: "pg"},
			Client:   clientMock,
		}
		targetCollection = "copy"
	)

	clientMock.EXPECT().
		WriteCollectionToResource(gomock.Any(), &pb.WriteCollectionRequest{
			Resource: r.Resource,
			SourceCollection: &pb.Collection{
				Name:    "name",
				Stream:  "stream",
				Records: []*pb.Record{},
			},
			TargetCollection: targetCollection,
			Configs: &pb.Configs{
				Config: []*pb.Config{},
			},
		}).Times(1).
		Return(&emptypb.Empty{}, nil)

	require.NoError(t, r.Write(sdk.Records{
		Name:    "name",
		Stream:  "stream",
		Records: []sdk.Record{},
	}, targetCollection))
}
