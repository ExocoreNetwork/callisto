package epochs

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/callisto/v4/database"

	epochssource "github.com/forbole/callisto/v4/modules/epochs/source"
	"github.com/forbole/juno/v5/modules"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.BlockModule   = &Module{}
)

// Module implements x/epochs module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source epochssource.Source
}

// NewModule builds a new Module instance
func NewModule(source epochssource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "epochs"
}
