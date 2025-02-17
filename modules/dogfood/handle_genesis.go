package dogfood

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	keytypes "github.com/ExocoreNetwork/exocore/types/keys"
	dogfoodtypes "github.com/ExocoreNetwork/exocore/x/dogfood/types"
	callistotypes "github.com/forbole/callisto/v4/types"
	junotypes "github.com/forbole/juno/v5/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", m.Name()).Msg("parsing genesis")

	// Read the genesis state
	var genState dogfoodtypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[dogfoodtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading dogfood genesis data: %s", err)
	}

	// Save the genesis params
	err = m.db.SaveDogfoodParams(&genState.Params, doc.InitialHeight)
	if err != nil {
		return fmt.Errorf("error while storing genesis dogfood params: %s", err)
	}

	// Initial validator vote powers (validators without the power should be
	// loaded from Tendermint API).
	// Since it is not clear if HandleGenesis is called after the validators
	// are loaded from Tendermint API, we add the validators (minus the vote power)
	// to the database first, and then, add the vote power separately.
	validators := make([]*junotypes.Validator, len(genState.ValSet))
	votingPowers := make([]callistotypes.ValidatorVotingPower, len(validators))
	for i, validator := range genState.ValSet {
		wrappedKey := keytypes.NewWrappedConsKeyFromHex(validator.PublicKey)
		if wrappedKey == nil {
			return fmt.Errorf("error while creating wrapped key: %s", validator.PublicKey)
		}
		consPubKey, err := junotypes.ConvertValidatorPubKeyToBech32String(wrappedKey.ToTmKey())
		if err != nil {
			return fmt.Errorf("error while converting validator pubkey to bech32 string: %s", err)
		}
		consAddr := wrappedKey.ToConsAddr().String()
		validators[i] = junotypes.NewValidator(
			consAddr, consPubKey,
		)
		votingPowers[i] = callistotypes.NewValidatorVotingPower(
			consAddr,
			sdk.TokensFromConsensusPower(validator.Power, sdk.DefaultPowerReduction).Int64(),
			doc.InitialHeight,
		)
	}
	// this function is the original one from x/staking or consensus
	// it never deletes a validator; instead, to check if a validator is currently active,
	// refer to its voting power at the most recent height on which it is available.
	err = m.db.SaveValidators(validators)
	if err != nil {
		return fmt.Errorf("error while storing genesis dogfood validators: %s", err)
	}

	// then we do the vote powers using the same x/staking function
	err = m.db.SaveValidatorsVotingPowers(votingPowers)
	if err != nil {
		return fmt.Errorf("error while storing genesis dogfood validators voting powers: %s", err)
	}

	// Opt out expiries
	for _, levelOne := range genState.OptOutExpiries {
		for _, levelTwo := range levelOne.OperatorAccAddrs {
			wrapped := callistotypes.NewOptOutExpiryFromStr(
				fmt.Sprintf("%d", levelOne.Epoch), levelTwo,
			)
			err = m.db.SaveOptOutExpiry(wrapped)
			if err != nil {
				return fmt.Errorf("error while storing genesis dogfood opt out expiry: %s", err)
			}
		}
	}

	// Consensus addresses to prune
	for _, levelOne := range genState.ConsensusAddrsToPrune {
		for _, levelTwo := range levelOne.ConsAddrs {
			wrapped := callistotypes.NewConsensusAddrToPruneFromStr(
				fmt.Sprintf("%d", levelOne.Epoch), levelTwo,
			)
			err = m.db.SaveConsensusAddrToPrune(wrapped)
			if err != nil {
				return fmt.Errorf("error while storing genesis dogfood consensus addr to prune: %s", err)
			}
		}
	}

	// Undelegation maturities
	for _, levelOne := range genState.UndelegationMaturities {
		for _, levelTwo := range levelOne.UndelegationRecordKeys {
			wrapped := callistotypes.NewUndelegationMaturityFromStr(
				fmt.Sprintf("%d", levelOne.Epoch), levelTwo,
			)
			err = m.db.SaveUndelegationMaturity(wrapped)
			if err != nil {
				return fmt.Errorf("error while storing genesis dogfood undelegation maturity: %s", err)
			}
		}
	}

	// Last total power
	err = m.db.SaveLastTotalPower(genState.LastTotalPower.String())
	if err != nil {
		return fmt.Errorf("error while storing genesis dogfood last total power: %s", err)
	}

	return nil
}
