package v3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v10/packetforward/types"

	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
)

// Migrate migrates the x/packetforward module state from the consensus version
// 2 to version 3
func Migrate(
	ctx sdk.Context,
	bankKeeper types.BankKeeper,
	channelKeeper types.ChannelKeeper,
	transferKeeper types.TransferKeeper,
) error {
	logger := ctx.Logger()

	expectedTotalEscrowed := sdk.Coins{}

	// 1. Iterate over all IBC transfer channels
	portID := transferKeeper.GetPort(ctx)
	transferChannels := channelKeeper.GetAllChannelsWithPortPrefix(ctx, portID)
	for _, channel := range transferChannels {
		// 2. For each channel, get the escrow address and corresponding bank balance
		escrowAddress := transfertypes.GetEscrowAddress(portID, channel.ChannelId)
		bankBalances := bankKeeper.GetAllBalances(ctx, escrowAddress)

		// 3. Aggregate the bank balances to calculate the expected total escrowed
		expectedTotalEscrowed = expectedTotalEscrowed.Add(bankBalances...)
	}

	logger.Info(
		"Calculated expected total escrowed from escrow account bank balances",
		"num channels", len(transferChannels),
		"bank total escrowed", expectedTotalEscrowed,
	)

	// 4. Set the total escrowed for each denom
	for _, totalEscrowCoin := range expectedTotalEscrowed {
		prevDenomEscrow := transferKeeper.GetTotalEscrowForDenom(ctx, totalEscrowCoin.Denom)

		transferKeeper.SetTotalEscrowForDenom(ctx, totalEscrowCoin)

		logger.Info(
			"Corrected total escrow for denom to match escrow account bank balances",
			"denom", totalEscrowCoin.Denom,
			"previous escrow", prevDenomEscrow,
			"new escrow", totalEscrowCoin,
		)
	}

	return nil
}
