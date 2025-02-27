package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	v3 "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v10/packetforward/migrations/v3"
)

// Migrator is a struct for handling in-place state migrations.
type Migrator struct {
	keeper *Keeper
}

func NewMigrator(k *Keeper) Migrator {
	return Migrator{
		keeper: k,
	}
}

// Migrate2to3 migrates the module state from the consensus version 2 to
// version 3
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v3.Migrate(
		ctx,
		m.keeper.bankKeeper,
		m.keeper.channelKeeper,
		m.keeper.transferKeeper,
	)
}
