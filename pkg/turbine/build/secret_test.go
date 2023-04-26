package build

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/client/mock"
)

func TestRegisterSecret(t *testing.T) {
	var (
		ctrl       = gomock.NewController(t)
		clientMock = mock.NewMockClient(ctrl)

		envName = "REGISTER_SECRET_TEST"
		envVal  = "some value"
	)

	os.Setenv(envName, envVal)
	defer os.Unsetenv(envName)

	clientMock.EXPECT().
		RegisterSecret(gomock.Any(), &pb.Secret{
			Name:  envName,
			Value: envVal,
		}).Times(1).
		Return(&emptypb.Empty{}, nil)
	b := builder{Client: clientMock}
	require.NoError(t, b.RegisterSecret(envName))
}
