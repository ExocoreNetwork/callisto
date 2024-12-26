package epochs

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"

	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"

	"github.com/rs/zerolog/log"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", m.Name()).Msg("parsing genesis")

	// Read the genesis state
	var genState epochstypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[epochstypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while unmarshalling staking state: %s", err)
	}

	// Save the epochs
	err = m.db.SaveEpochs(genState.Epochs)
	if err != nil {
		return fmt.Errorf("error while storing genesis staking params: %s", err)
	}

	return nil
}
