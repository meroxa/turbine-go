package turbine

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/turbine-go/pkg/app"
	"github.com/meroxa/turbine-go/pkg/proto/core"
	"github.com/meroxa/turbine-go/pkg/turbine/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestResources(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		turbineMock = mock.NewMockTurbineCore(ctrl)

		resourceName = "pg"
		coreResource = &core.Resource{
			Name: resourceName,
		}
	)

	turbineMock.EXPECT().
		GetResource(gomock.Any(), &core.GetResourceRequest{
			Name: resourceName,
		}).Times(1).
		Return(coreResource, nil)

	tb := turbine{TurbineCore: turbineMock}
	r, err := tb.Resources(resourceName)
	require.NoError(t, err)
	require.Equal(t, r, &resource{
		Resource: coreResource,
		tc:       turbineMock,
	})
}

func TestRecords(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		turbineMock = mock.NewMockTurbineCore(ctrl)

		r = &resource{
			Resource: &core.Resource{
				Name: "pg",
			},
			tc: turbineMock,
		}
		collection = "accounts"
	)

	turbineMock.EXPECT().
		ReadCollection(gomock.Any(), &core.ReadCollectionRequest{
			Resource:   r.Resource,
			Collection: collection,
			Configs: &core.Configs{
				Config: []*core.Config{
					{
						Field: "conf",
						Value: "conf_val",
					},
				},
			},
		}).Times(1).
		Return(&core.Collection{
			Name:   "name",
			Stream: "stream",
		}, nil)

	c, err := r.Records(collection, app.ConnectionOptions{{
		Field: "conf",
		Value: "conf_val",
	}})
	require.NoError(t, err)
	require.Equal(t, c, app.Records{
		Name:    "name",
		Stream:  "stream",
		Records: []app.Record{},
	})
}

func TestWrite(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		turbineMock = mock.NewMockTurbineCore(ctrl)

		r = &resource{
			Resource: &core.Resource{
				Name: "pg",
			},
			tc: turbineMock,
		}
		targetCollection = "copy"
	)

	turbineMock.EXPECT().
		WriteCollectionToResource(gomock.Any(), &core.WriteCollectionRequest{
			Resource: r.Resource,
			SourceCollection: &core.Collection{
				Name:    "name",
				Stream:  "stream",
				Records: []*core.Record{},
			},
			TargetCollection: targetCollection,
			Configs: &core.Configs{
				Config: []*core.Config{},
			},
		}).Times(1).
		Return(&emptypb.Empty{}, nil)

	err := r.Write(app.Records{
		Name:    "name",
		Stream:  "stream",
		Records: []app.Record{},
	}, targetCollection)
	require.NoError(t, err)
}
