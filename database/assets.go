package database

import (
	"encoding/json"
	"fmt"

	sdkmath "cosmossdk.io/math"

	"github.com/forbole/callisto/v4/types"
)

// SaveAssetsParams allows to store the given params inside the database
func (db *Db) SaveAssetsParams(params *types.AssetsParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling assets params: %s", err)
	}

	stmt := `
INSERT INTO assets_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE assets_params.height <= excluded.height`

	_, err = db.SQL.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing assets params: %s", err)
	}

	return nil
}

// SaveClientChain inserts or updates a client chain record in the database
func (db *Db) SaveOrUpdateClientChain(chain *types.ClientChain) error {
	stmt := `
INSERT INTO client_chains (name, meta_info, chain_id, exocore_chain_index, finalization_blocks, layer_zero_chain_id, signature_type, address_length)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (layer_zero_chain_id) DO UPDATE
	SET name = EXCLUDED.name,
		meta_info = EXCLUDED.meta_info,
		chain_id = EXCLUDED.chain_id,
		finalization_blocks = EXCLUDED.finalization_blocks,
		layer_zero_chain_id = EXCLUDED.layer_zero_chain_id,
		signature_type = EXCLUDED.signature_type,
		address_length = EXCLUDED.address_length;`

	_, err := db.SQL.Exec(stmt,
		chain.Name,
		chain.MetaInfo,
		chain.ChainId,
		chain.ExocoreChainIndex,
		chain.FinalizationBlocks,
		chain.LayerZeroChainID,
		chain.SignatureType,
		chain.AddressLength,
	)
	if err != nil {
		return fmt.Errorf("failed to save client chain: %w", err)
	}
	return nil
}

// SaveAssetsToken saves a token record into the database. Once added, only the
// metadata may be altered by the blockchain.
func (db *Db) SaveAssetsToken(token *types.AssetsToken) error {
	// Q. drop the total deposit amount or retain?
	// A. slashing is applied to staker level, not the deposit amount.
	//    so it is a good thing to retain the total deposit amount.
	stmt := `
INSERT INTO assets_tokens (asset_id, name, symbol, address, decimals, layer_zero_chain_id, exocore_chain_index, meta_info, staking_total_amount)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (asset_id) DO UPDATE
SET name = EXCLUDED.name,
    symbol = EXCLUDED.symbol,
    address = EXCLUDED.address,
    decimals = EXCLUDED.decimals,
    layer_zero_chain_id = EXCLUDED.layer_zero_chain_id,
    exocore_chain_index = EXCLUDED.exocore_chain_index,
    meta_info = EXCLUDED.meta_info
	staking_total_amount = EXCLUDED.staking_total_amount;`
	_, err := db.SQL.Exec(stmt,
		token.AssetID,
		token.Name,
		token.Symbol,
		token.Address,
		token.Decimals,
		token.LayerZeroChainID,
		token.ExocoreChainIndex,
		token.MetaInfo,
		token.Amount,
	)
	if err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}
	return nil
}

// UpdateAssetMetadata updates the metadata for an asset based on its assetID.
func (db *Db) UpdateAssetMetadata(
	assetID string, newMetaInfo string,
) error {
	stmt := `
UPDATE tokens
SET meta_info = $1
WHERE asset_id = $2;`
	// no conflict can happen above since it is update
	// we don't check for existence because again, we are not a business logic implementation
	_, err := db.SQL.Exec(stmt, newMetaInfo, assetID)
	if err != nil {
		return fmt.Errorf(
			"failed to update metadata for asset_id %s: %w",
			assetID, err,
		)
	}
	return nil
}

// UpdateStakingTotalAmount updates the total staking amount for an asset based on its assetID.
func (db *Db) UpdateStakingTotalAmount(
	assetID string, newAmount string,
) error {
	stmt := `
UPDATE tokens
SET staking_total_amount = $1
WHERE asset_id = $2;`
	_, err := db.SQL.Exec(stmt, newAmount, assetID)
	if err != nil {
		return fmt.Errorf(
			"failed to update staking total amount for asset_id %s: %w",
			assetID, err,
		)
	}
	return nil
}

// SaveStakerAsset saves a staker asset record in the database, including
// an entry in the history table.
func (db *Db) SaveStakerAsset(data *types.StakerAsset) error {
	stmt := `
INSERT INTO staker_assets (staker_id, asset_id, deposited, withdrawable, pending_undelegation)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (staker_id, asset_id) DO UPDATE
SET deposited = EXCLUDED.deposited,
	withdrawable = EXCLUDED.withdrawable,
	pending_undelegation = EXCLUDED.pending_undelegation,
	delegated = EXCLUDED.deposited - EXCLUDED.withdrawable - EXCLUDED.pending_undelegation - staker_assets.lifetime_slashed,
	lifetime_slashed = staker_assets.lifetime_slashed;`

	_, err := db.SQL.Exec(
		stmt,
		data.StakerID,
		data.AssetID,
		data.Deposited,
		data.Withdrawable,
		data.PendingUndelegation,
	)
	if err != nil {
		return fmt.Errorf("failed to save staker asset: %w", err)
	}
	return nil
}

// GetDelegatedAmount returns the delegated amount for a given staker and asset.
func (db *Db) GetDelegatedAmount(stakerID, assetID string) (sdkmath.Int, error) {
	stmt := `
	SELECT delegated FROM staker_assets WHERE staker_id = $1 AND asset_id = $2;`
	var delegatedAmount string
	err := db.SQL.QueryRow(stmt, stakerID, assetID).Scan(&delegatedAmount)
	if err != nil {
		return sdkmath.Int{}, fmt.Errorf("failed to get delegated amount: %w", err)
	}
	delegatedAmountInt, ok := sdkmath.NewIntFromString(delegatedAmount)
	if !ok {
		return sdkmath.Int{}, fmt.Errorf("failed to convert delegated amount to int: %s", delegatedAmount)
	}
	return delegatedAmountInt, nil
}

// SlashStakerDelegation slashes the staker delegation. It updates the lifetime slashed amount
// and the delegated amount.
func (db *Db) SlashStakerDelegation(stakerID, assetID, slashedAmount string) error {
	stmt := `
	UPDATE staker_assets
	SET lifetime_slashed = lifetime_slashed + $1,
	    delegated = delegated - $1
	WHERE staker_id = $2 AND asset_id = $3;`
	_, err := db.SQL.Exec(stmt, slashedAmount, stakerID, assetID)
	if err != nil {
		return fmt.Errorf("failed to accumulate staker lifetime slashing: %w", err)
	}
	return nil
}

// SaveOperatorAsset saves an operator asset record in the database,
// ensuring other_share is derived as total_share - self_share.
// This function is in `assets.go` because it is triggered by events in `x/assets`.
func (db *Db) SaveOperatorAsset(data *types.OperatorAsset) error {
	stmt := `
INSERT INTO operator_assets (operator_addr, asset_id, total_amount, pending_undelegation_amount, total_share, self_share, other_share)
VALUES ($1, $2, $3, $4, $5, $6, $5 - $6)
ON CONFLICT (operator_addr, asset_id) DO UPDATE
SET total_amount = EXCLUDED.total_amount,
	pending_undelegation_amount = EXCLUDED.pending_undelegation_amount,
	total_share = EXCLUDED.total_share,
	self_share = EXCLUDED.self_share,
	other_share = EXCLUDED.total_share - EXCLUDED.self_share;`
	_, err := db.SQL.Exec(
		stmt, data.OperatorAddress, data.AssetID,
		data.TotalAmount, data.PendingUndelegationAmount,
		data.TotalShare, data.SelfShare,
	)
	if err != nil {
		return fmt.Errorf("failed to save operator asset: %w", err)
	}
	return nil
}
