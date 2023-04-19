package turbine

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/turbine-go/pkg/proto/core"
	"github.com/meroxa/turbine-go/pkg/turbine/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestRegisterSecret(t *testing.T) {
	var (
		ctrl        = gomock.NewController(t)
		turbineMock = mock.NewMockTurbineCore(ctrl)

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

	tb := turbine{TurbineCore: turbineMock}
	err := tb.RegisterSecret(envName)
	require.NoError(t, err)
}
