package dogfood

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/callisto/v4/database"

	dogfoodsource "github.com/forbole/callisto/v4/modules/dogfood/source"
	"github.com/forbole/juno/v5/modules"
)

var (
	_ modules.Module            = &Module{}
	_ modules.GenesisModule     = &Module{}
	_ modules.TransactionModule = &Module{}
	_ modules.BlockModule       = &Module{}
	_ modules.MessageModule     = &Module{}
)

// Module implements x/assets module indexer
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source dogfoodsource.Source
}

// NeawModule builds a new Module instance
func NewModule(source dogfoodsource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "dogfood"
}
