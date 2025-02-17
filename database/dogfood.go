package database

import (
	"fmt"

	"github.com/lib/pq"

	"github.com/forbole/callisto/v4/types"

	dogfoodtypes "github.com/ExocoreNetwork/exocore/x/dogfood/types"
)

// SaveDogfoodParams saves the dogfood params. It is called upon genesis and post
// a governance vote pass to change the params.
func (db *Db) SaveDogfoodParams(params *dogfoodtypes.Params, height int64) error {
	stmt := `
	INSERT INTO dogfood_params (one_row_id, height, epochs_until_unbonded, epoch_identifier, max_validators, historical_entries, min_self_delegation, asset_ids)
	VALUES (TRUE, $1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (one_row_id) DO UPDATE
	SET height = $1,
	    epochs_until_unbonded = $2,
	    epoch_identifier = $3,
	    max_validators = $4,
	    historical_entries = $5,
	    min_self_delegation = $6,
	    asset_ids = $7;`
	_, err := db.SQL.Exec(
		stmt,
		height,
		params.EpochsUntilUnbonded,
		params.EpochIdentifier,
		params.MaxValidators,
		params.HistoricalEntries,
		params.MinSelfDelegation.String(),
		pq.Array(params.AssetIDs),
	)
	if err != nil {
		return fmt.Errorf("failed to save dogfood params: %w", err)
	}
	return nil
}

// SaveOptOutExpiry saves the opt out expiry for a given operator.
func (db *Db) SaveOptOutExpiry(optOut *types.OptOutExpiry) error {
	// for each operator, we can have only one opt out expiry active.
	// ideally there should never be a conflict, but that is left to the chain.
	stmt := `
	INSERT INTO opt_out_expiries (epoch_number, operator_addr)
	VALUES ($1, $2)
	ON CONFLICT (operator_addr) DO UPDATE
	SET epoch_number = $1;`
	_, err := db.SQL.Exec(
		stmt,
		optOut.EpochNumber,
		optOut.OperatorAddr,
	)
	if err != nil {
		return fmt.Errorf("failed to save opt out expiry: %w", err)
	}
	return nil
}

// CompleteOptOuts completes the opt out for a given epoch.
func (db *Db) CompleteOptOuts(epochNumber string, height int64) error {
	stmt := `
	UPDATE opt_out_expiries
	SET completion_height = $1
	WHERE epoch_number = $2;`
	_, err := db.SQL.Exec(
		stmt,
		height, epochNumber,
	)
	if err != nil {
		return fmt.Errorf("failed to complete opt outs: %w", err)
	}
	return nil
}

// SaveConsensusAddrToPrune saves the consensus address to prune for a given epoch.
func (db *Db) SaveConsensusAddrToPrune(consAddr *types.ConsensusAddrToPrune) error {
	stmt := `
	INSERT INTO consensus_addrs_to_prune (epoch_number, consensus_addr)
	VALUES ($1, $2)
	ON CONFLICT (consensus_addr) DO UPDATE
	SET epoch_number = $1;`
	_, err := db.SQL.Exec(
		stmt,
		consAddr.EpochNumber,
		consAddr.ConsensusAddr,
	)
	if err != nil {
		return fmt.Errorf("failed to save consensus addr to prune: %w", err)
	}
	return nil
}

// CompleteConsensusAddrPruning completes the consensus address pruning for a given epoch.
func (db *Db) CompleteConsensusAddrsPruning(epochNumber string, height int64) error {
	stmt := `
	UPDATE consensus_addrs_to_prune
	SET completion_height = $1
	WHERE epoch_number = $2;`
	_, err := db.SQL.Exec(
		stmt,
		height, epochNumber,
	)
	if err != nil {
		return fmt.Errorf("failed to complete consensus addr pruning: %w", err)
	}
	return nil
}

// SaveUndelegationMaturity saves the undelegation maturity for a given operator.
func (db *Db) SaveUndelegationMaturity(undelegation *types.UndelegationMaturity) error {
	stmt := `
	INSERT INTO undelegation_maturities (epoch_number, record_key)
	VALUES ($1, $2)
	ON CONFLICT (record_key) DO UPDATE
	SET epoch_number = $1;`
	_, err := db.SQL.Exec(
		stmt,
		undelegation.EpochNumber,
		undelegation.RecordKey,
	)
	if err != nil {
		return fmt.Errorf("failed to save undelegation maturity: %w", err)
	}
	return nil
}

// MatureUndelegations matures the undelegation for a given record key.
func (db *Db) MatureUndelegations(epochNumber string, height int64) error {
	stmt := `
	UPDATE undelegation_maturities
	SET completion_height = $1
	WHERE epoch_number = $2;`
	_, err := db.SQL.Exec(
		stmt,
		height, epochNumber,
	)
	if err != nil {
		return fmt.Errorf("failed to mature undelegations: %w", err)
	}
	return nil
}

// SaveLastTotalPower saves the last total power.
func (db *Db) SaveLastTotalPower(lastTotalPower string) error {
	stmt := `
	INSERT INTO last_total_power (one_row_id, total_power)
	VALUES (TRUE, $1)
	ON CONFLICT (one_row_id) DO UPDATE
	SET total_power = $1;`
	_, err := db.SQL.Exec(
		stmt,
		lastTotalPower,
	)
	if err != nil {
		return fmt.Errorf("failed to save last total power: %w", err)
	}
	return nil
}
