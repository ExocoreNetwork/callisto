package exomint

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/forbole/callisto/v4/database"

	"github.com/forbole/juno/v5/modules"

	exomintsource "github.com/forbole/callisto/v4/modules/exomint/source"
)

var (
	_ modules.Module             = &Module{}
	_ modules.GenesisModule      = &Module{}
	_ modules.BlockModule        = &Module{}
	_ modules.MessageModule      = &Module{}
	_ modules.AuthzMessageModule = &Module{}
)

// Module implements x/exomint module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source exomintsource.Source
}

// NewModule builds a new Module instance
func NewModule(source exomintsource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "exomint"
}
