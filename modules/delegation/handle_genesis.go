package delegation

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/forbole/callisto/v4/types"

	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", m.Name()).Msg("parsing genesis")

	// Read the genesis state
	var genState delegationtypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[delegationtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading delegation genesis data: %s", err)
	}

	// associations
	for _, association := range genState.Associations {
		if err := m.db.SaveStakerOperatorAssociation(association.StakerID, association.Operator); err != nil {
			return fmt.Errorf("error while saving association: %s", err)
		}
	}

	// delegation states
	for _, state := range genState.DelegationStates {
		keys, err := delegationtypes.ParseStakerAssetIDAndOperator([]byte(state.Key))
		if err != nil {
			return fmt.Errorf("error while parsing delegation state key: %s", err)
		}
		wrappedState := types.NewDelegationState(keys.StakerID, keys.AssetID, keys.OperatorAddr, &state.States)
		if err := m.db.SaveDelegationState(wrappedState); err != nil {
			return fmt.Errorf("error while saving delegation state: %s", err)
		}
		// capture the exoAssetDelegation state manually.
		// this is because the exoAssetDelegation is not tracked in the assets module.
		// the delegation module tracks it with a different logic, which we work around here.
		if keys.AssetID == assetstypes.ExocoreAssetID {
			// TODO: instead of GetDelegatedAmount, can we consider `TotalDelegatedAmountForStakerAsset` ?
			// the advantage of that function, if implemented, is that it does not need the operator address.
			delegatedAmount, err := m.source.GetDelegatedAmount(
				doc.InitialHeight, keys.StakerID, keys.AssetID, keys.OperatorAddr,
			)
			if err != nil {
				return fmt.Errorf("error while getting delegated amount: %s", err)
			}
			delegation := types.NewExoAssetDelegationFromStr(
				keys.StakerID, delegatedAmount.String(),
				state.States.WaitUndelegationAmount.String(),
			)
			if err := m.db.AccumulateExoAssetDelegation(delegation); err != nil {
				return fmt.Errorf("error while accumulating exo asset delegation: %s", err)
			}
		}
	}

	// stakers for each operator
	for _, data := range genState.StakersByOperator {
		parsed, err := assetstypes.ParseJoinedStoreKey([]byte(data.Key), 2)
		if err != nil {
			return fmt.Errorf("error while parsing staker by operator key: %s", err)
		}
		operatorAddr, assetID := parsed[0], parsed[1]
		for _, stakerID := range data.Stakers {
			if err := m.db.AppendStakerToOperatorAsset(stakerID, operatorAddr, assetID); err != nil {
				return fmt.Errorf("error while appending staker to operator asset: %s", err)
			}
		}
	}

	// undelegation records and hold counts
	for _, record := range genState.Undelegations {
		// in v1.0.9, the hold count is not part of the genesis state.
		wrapped := types.NewUndelegationRecord(&record, 0)
		if err := m.db.SaveUndelegationRecord(wrapped); err != nil {
			return fmt.Errorf("error while saving undelegation record: %s", err)
		}
	}

	// last known undelegation id, probably skip this one

	return nil
}
