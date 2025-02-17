package avs

import (
	"encoding/json"
	"fmt"
	"strings"

	tmtypes "github.com/cometbft/cometbft/types"

	avstypes "github.com/ExocoreNetwork/exocore/x/avs/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", m.Name()).Msg("parsing genesis")
	var state avstypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[avstypes.ModuleName], &state)
	if err != nil {
		return fmt.Errorf("error while reading avs genesis data: %s", err)
	}
	// TODO: handle this completely
	for _, avs := range state.AvsInfos {
		avs.AvsAddress = strings.ToLower(avs.AvsAddress)
		m.db.SaveAvsAddr(avs.AvsAddress)
	}
	for _, elem := range state.ChainIdInfos {
		elem.AvsAddress = strings.ToLower(elem.AvsAddress)
		m.db.SaveChainIdToAvsAddr(elem.ChainId, elem.AvsAddress)
	}
	return nil
}
