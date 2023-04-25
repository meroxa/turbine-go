package turbine

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/client/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestRegisterSecret(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		turbineMock = mock.NewMockClient(ctrl)

		envName = "REGISTER_SECRET_TEST"
		envVal  = "some value"
	)

	os.Setenv(envName, envVal)
	defer os.Unsetenv(envName)

	turbineMock.EXPECT().
		RegisterSecret(gomock.Any(), &core.Secret{
			Name:  envName,
			Value: envVal,
		}).Times(1).
		Return(&emptypb.Empty{}, nil)

	tb := turbine{Client: turbineMock}
	err := tb.RegisterSecret(envName)
	require.NoError(t, err)
}
