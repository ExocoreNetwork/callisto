package avs

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/callisto/v4/database"

	"github.com/forbole/juno/v5/modules"
)

var (
	_ modules.Module            = &Module{}
	_ modules.GenesisModule     = &Module{}
	_ modules.TransactionModule = &Module{}
	_ modules.BlockModule       = &Module{}
)

// Module implements x/assets module indexer
type Module struct {
	cdc codec.Codec
	db  *database.Db
}

// NewModule builds a new Module instance
func NewModule(cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc: cdc,
		db:  db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "avs"
}
