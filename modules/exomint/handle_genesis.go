package exomint

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/forbole/callisto/v4/types"

	exominttypes "github.com/ExocoreNetwork/exocore/x/exomint/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", m.Name()).Msg("parsing genesis")

	// Read the genesis state
	var genState exominttypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[exominttypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading exomint genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveExomintParams(types.NewExomintParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis exomint params: %s", err)
	}

	return nil
}
