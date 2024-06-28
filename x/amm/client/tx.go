package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/amm/types"
	"github.com/spf13/cobra"
	"strconv"
)

var (
	FlagWithdrawAll = "withdrawAll"
)

func GetTxCmd() *cobra.Command {
	AuthorizationTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "AMM module transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	AuthorizationTxCmd.AddCommand(
		NewCmdCreatePool(),
		NewCmdJoinPool(),
		NewCmdSwap(),
		NewCmdExitPool(),
	)

	return AuthorizationTxCmd
}

func NewCmdCreatePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [token_1] [token_2] [fee]",
		Short: "create new pool",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			token1, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			token2, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			fee, err := sdk.NewDecFromStr(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgCreatePoolMessage(clientCtx.GetFromAddress(), token1, token2, fee)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewCmdJoinPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "join [poolId] [token]",
		Short: "join an existing pool",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			token, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgJoinPoolMessage(clientCtx.GetFromAddress(), poolId, token)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewCmdSwap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [poolId] [token]",
		Short: "swap in a pool",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			token, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgSwapMessage(clientCtx.GetFromAddress(), poolId, token)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewCmdExitPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exit [poolId] [lpShare]",
		Short: "exit a pool",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			lpShare, err := sdk.NewDecFromStr(args[1])
			if err != nil {
				return err
			}

			withdrawAll, _ := cmd.Flags().GetBool(FlagWithdrawAll)

			msg := types.NewMsgExitPoolMessage(clientCtx.GetFromAddress(), poolId, lpShare, withdrawAll)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().Bool(FlagWithdrawAll, false, "Withdraw all from pool")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
