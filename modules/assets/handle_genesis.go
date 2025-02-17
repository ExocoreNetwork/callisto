package assets

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/forbole/callisto/v4/types"

	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", m.Name()).Msg("parsing genesis")

	// Read the genesis state
	var genState assetstypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[assetstypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading assets genesis data: %s", err)
	}

	// the genesis is made up of
	// - params
	// - client chains
	// - tokens
	// - deposits
	// - operator states
	// while regularly, these states are handled using events, genesis events are
	// discarded by CometBFT, so we must handle the genesis state here.

	// Save the genesis params
	err = m.db.SaveAssetsParams(types.NewAssetsParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis assets params: %s", err)
	}

	// Save the genesis client chains
	for _, chain := range genState.ClientChains {
		err = m.db.SaveOrUpdateClientChain(types.NewClientChain(chain))
		if err != nil {
			return fmt.Errorf("error while storing genesis client chain: %s", err)
		}
	}

	// Save the genesis tokens
	for _, token := range genState.Tokens {
		err = m.db.SaveAssetsToken(types.NewAssetsToken(&token))
		if err != nil {
			return fmt.Errorf("error while storing genesis token: %s", err)
		}
	}

	// Save the genesis deposits
	for _, layerOne := range genState.Deposits {
		stakerID := layerOne.StakerID
		for _, layerTwo := range layerOne.Deposits {
			assetID := layerTwo.AssetID
			deposit := layerTwo.Info
			err = m.db.SaveStakerAsset(
				types.NewStakerAssetFromInfo(
					stakerID, assetID, deposit,
				),
			)
			if err != nil {
				return fmt.Errorf("error while storing genesis deposit: %s", err)
			}
		}
	}

	// Save the genesis operator states
	for _, layerOne := range genState.OperatorAssets {
		operatorAddress := layerOne.Operator
		for _, layerTwo := range layerOne.AssetsState {
			assetID := layerTwo.AssetID
			state := layerTwo.Info
			err = m.db.SaveOperatorAsset(
				types.NewOperatorAssetFromInfo(
					operatorAddress, assetID, state,
				),
			)
			if err != nil {
				return fmt.Errorf("error while storing genesis operator state: %s", err)
			}
		}
	}

	return nil
}
