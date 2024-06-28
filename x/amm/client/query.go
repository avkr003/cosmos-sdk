package client

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/amm/types"
	"github.com/spf13/cobra"
	"strconv"
)

func GetQueryCmd() *cobra.Command {
	authorizationQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the amm module",
		Long:                       "",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	authorizationQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryPool(),
		GetCmdQueryPoolShares(),
	)

	return authorizationQueryCmd
}

func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the amm module parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Parameters(cmd.Context(), &types.QueryParametersRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool [poolId]",
		Args:  cobra.ExactArgs(1),
		Short: "Query a pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			poolId, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil || poolId < 0 {
				return fmt.Errorf("poolId argument provided must be a non-negative-integer: %v", err)
			}

			res, err := queryClient.Pool(cmd.Context(), &types.QueryPoolRequest{PoolId: uint64(poolId)})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Pool)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryPoolShares() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "poolShares [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Query pool shares for address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			res, err := queryClient.PoolShares(cmd.Context(), &types.QueryPoolSharesRequest{Address: address.String(), Pagination: pageReq})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all pool shares")

	return cmd
}
