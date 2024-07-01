package keeper_test

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/amm/keeper"
	ammtestutil "github.com/cosmos/cosmos-sdk/x/amm/testutil"
	"github.com/cosmos/cosmos-sdk/x/amm/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	ammKeeper   keeper.Keeper
	bankKeeper  *ammtestutil.MockBankKeeper
	queryClient types.QueryClient
}

func (s *KeeperTestSuite) SetupTest() {
	key := sdk.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, sdk.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()

	ctrl := gomock.NewController(s.T())
	bankKeeper := ammtestutil.NewMockBankKeeper(ctrl)

	keeper := keeper.NewKeeper(
		encCfg.Codec,
		key,
		bankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	keeper.SetParams(ctx, types.DefaultParams())

	s.ctx = ctx
	s.bankKeeper = bankKeeper
	s.ammKeeper = keeper

	types.RegisterInterfaces(encCfg.InterfaceRegistry)
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, keeper)
	s.queryClient = types.NewQueryClient(queryHelper)
}

func (s *KeeperTestSuite) TestParams() {
	ctx, ammKeeper := s.ctx, s.ammKeeper
	require := s.Require()

	expParams := types.DefaultParams()
	expParams.SwapAllowedTokens = []string{"ubtc", "ueth"}
	expParams.MaxSwapFee = sdk.MustNewDecFromStr("0.25")
	ammKeeper.SetParams(ctx, expParams)
	resParams := ammKeeper.GetParams(ctx)
	require.True(expParams.Equal(resParams))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
