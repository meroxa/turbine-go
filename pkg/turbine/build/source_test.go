package build

// TODO: Implement

//func TestSource(t *testing.T) {
//	var (
//		ctrl       = gomock.NewController(t)
//		clientMock = mock.NewMockClient(ctrl)
//
//		resourceName = "pg"
//		coreSource   = &pb.Source{
//			Name: resourceName,
//		}
//	)
//
//	clientMock.EXPECT().
//		GetSource(gomock.Any(), &pb.GetSourceRequest{
//			Name: resourceName,
//		}).Times(1).
//		Return(coreSource, nil)
//
//	b := builder{Client: clientMock}
//	r, err := b.Source(resourceName)
//	require.NoError(t, err)
//	require.Equal(t, r, &source{
//		Source: coreSource,
//		Client: clientMock,
//	})
//}
//
//func TestSourceRecords(t *testing.T) {
//	var (
//		ctrl       = gomock.NewController(t)
//		clientMock = mock.NewMockClient(ctrl)
//
//		r = &source{
//			Source: &pb.Source{Name: "pg"},
//			Client: clientMock,
//		}
//		collection = "accounts"
//	)
//
//	clientMock.EXPECT().
//		ReadCollection(gomock.Any(), &pb.ReadCollectionRequest{
//			Source:     r.Source,
//			Collection: collection,
//			Configs: &pb.Configs{
//				Config: []*pb.Config{
//					{
//						Field: "conf",
//						Value: "conf_val",
//					},
//				},
//			},
//		}).Times(1).
//		Return(&pb.Collection{
//			Name:   "name",
//			Stream: "stream",
//		}, nil)
//
//	c, err := r.Read(collection, sdk.ConnectionOptions{{
//		Field: "conf",
//		Value: "conf_val",
//	}})
//	require.NoError(t, err)
//	require.Equal(t, c, sdk.Records{
//		Name:    "name",
//		Stream:  "stream",
//		Records: []sdk.Record{},
//	})
//}
//
//func TestWriteToDestination(t *testing.T) {
//	var (
//		ctrl       = gomock.NewController(t)
//		clientMock = mock.NewMockClient(ctrl)
//
//		r = &source{
//			Source: &pb.Source{Name: "pg"},
//			Client: clientMock,
//		}
//		targetCollection = "copy"
//	)
//
//	clientMock.EXPECT().
//		WriteCollectionToDestination(gomock.Any(), &pb.WriteCollectionRequest{
//			Destination: r.Destination,
//			SourceCollection: &pb.Collection{
//				Name:    "name",
//				Stream:  "stream",
//				Records: []*pb.Record{},
//			},
//			DestinationCollection: targetCollection,
//			Configs: &pb.Configs{
//				Config: []*pb.Config{},
//			},
//		}).Times(1).
//		Return(&emptypb.Empty{}, nil)
//
//	require.NoError(t, r.Write(sdk.Records{
//		Name:    "name",
//		Stream:  "stream",
//		Records: []sdk.Record{},
//	}, targetCollection))
//}
