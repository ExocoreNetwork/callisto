package operator

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/forbole/callisto/v4/types"

	keytypes "github.com/ExocoreNetwork/exocore/types/keys"
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", m.Name()).Msg("parsing genesis")
	var genState operatortypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[operatortypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading operator genesis data: %s", err)
	}
	// genesis state is made up of:
	// - operator
	for _, operatorDetail := range genState.Operators {
		if err := m.db.SaveOperatorDetail(types.NewOperator(&operatorDetail.OperatorInfo)); err != nil {
			return fmt.Errorf("error while saving operator detail: %s", err)
		}
	}
	// - operator consensus keys for chain ids
	for _, record := range genState.OperatorRecords {
		addr := record.OperatorAddress
		for _, detail := range record.Chains {
			// already validated by the chain
			wrappedKey := keytypes.NewWrappedConsKeyFromHex(detail.ConsensusKey)
			// to avoid using our own chain's bech32-prefix, use the hex representation
			// non-checksummed without "0x" prefix
			consAddress := hex.EncodeToString(wrappedKey.ToConsAddr().Bytes())
			if err := m.db.SaveOperatorConsKey(addr, detail.ChainID, detail.ConsensusKey, consAddress); err != nil {
				return fmt.Errorf("error while saving operator cons key: %s", err)
			}
		}
	}
	// - operator opted in state
	for _, data := range genState.OptStates {
		keys, err := assetstypes.ParseJoinedStoreKey([]byte(data.Key), 2)
		if err != nil {
			return fmt.Errorf("failed to parse joined key: %w", err)
		}
		operatorAddr, avsAddr := keys[0], keys[1]
		opted := types.NewOpted(
			operatorAddr, avsAddr, &data.OptInfo,
		)
		if err := m.db.SaveOptedState(opted); err != nil {
			return fmt.Errorf("error while saving operator opt state: %s", err)
		}
	}
	// - operator usd values
	for _, usdValue := range genState.OperatorUSDValues {
		parsed, err := assetstypes.ParseJoinedStoreKey([]byte(usdValue.Key), 2)
		if err != nil {
			return fmt.Errorf("error while parsing operator usd value: %s", err)
		}
		avsAddr, operatorAddr := parsed[0], parsed[1]
		operatorUSDValue := types.NewOperatorUSDValue(
			operatorAddr, avsAddr,
			&usdValue.OptedUSDValue,
		)
		if err := m.db.SaveOperatorUSDValue(operatorUSDValue); err != nil {
			return fmt.Errorf("error while saving operator usd value: %s", err)
		}
	}
	// - avs usd values
	for _, usdValue := range genState.AVSUSDValues {
		avsUsdValue := types.NewAvsUSDValueFromStr(usdValue.AVSAddr, usdValue.Value.String())
		if err := m.db.SaveAvsUSDValue(avsUsdValue); err != nil {
			return fmt.Errorf("error while saving avs usd value: %s", err)
		}
	}

	// - TODO slash states (skipped for now)

	// - prev consensus keys
	// TODO: is this even worth tracking?
	// for any given operator, it represents the key from which they will switch
	// to another within this epoch. it is useful, for example, if you want to
	// track the key to the operator and find signing details.
	for _, prev := range genState.PreConsKeys {
		parsed, err := assetstypes.ParseJoinedStoreKey([]byte(prev.Key), 2)
		if err != nil {
			return fmt.Errorf("error while parsing prev cons key: %s", err)
		}
		chainId, operatorAddr := parsed[0], parsed[1]
		wrappedKey := keytypes.NewWrappedConsKeyFromHex(prev.ConsensusKey)
		if err := m.db.SaveOperatorPrevConsKey(
			chainId, operatorAddr, wrappedKey.ToHex(), wrappedKey.ToConsAddr().String(),
		); err != nil {
			return fmt.Errorf("error while saving prev cons key: %s", err)
		}
	}
	// - consensus key removals
	for _, removal := range genState.OperatorKeyRemovals {
		parsed, err := assetstypes.ParseJoinedStoreKey([]byte(removal.Key), 2)
		if err != nil {
			return fmt.Errorf("error while parsing operator key removal: %s", err)
		}
		// NOTE the reversed order
		operatorAddr, chainId := parsed[0], parsed[1]
		if err := m.db.MarkOperatorKeyRemoval(chainId, operatorAddr); err != nil {
			return fmt.Errorf("error while removing operator cons key: %s", err)
		}
	}
	return nil
}
