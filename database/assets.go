package database

import (
	"encoding/json"
	"fmt"

	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
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
func (db *Db) SaveOrUpdateClientChain(chain assetstypes.ClientChainInfo) error {
	// this can not be edited, except in the case of a chain restart.
	// in that event, the index should be regenerated.
	stmt := `
INSERT INTO client_chains (name, meta_info, chain_id, exocore_chain_index, finalization_blocks, layer_zero_chain_id, signature_type, address_length)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (exocore_chain_index) DO NOTHING`

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

// SaveToken saves a token record into the database. Once added, only the
// metadata may be altered by the blockchain.
func (db *Db) SaveToken(tokenInput assetstypes.StakingAssetInfo) error {
	// drop the total deposit amount or retain?
	// slashing is applied to staker level, not the deposit amount.
	// so it is a good thing to retain the total deposit amount.
	token := tokenInput.AssetBasicInfo
	stmt := `
INSERT INTO tokens (asset_id, name, symbol, address, decimals, layer_zero_chain_id, exocore_chain_index, meta_info, staking_total_amount)
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
		token.AssetID(),
		token.Name,
		token.Symbol,
		token.Address,
		token.Decimals,
		token.LayerZeroChainID,
		token.ExocoreChainIndex,
		token.MetaInfo,
		tokenInput.StakingTotalAmount,
	)
	if err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}
	return nil
}

// UpdateAssetMetadata updates the metadata for an asset based on its
// address and LayerZero chain ID.
func (db *Db) UpdateAssetMetadata(
	address string, layerZeroChainID int64, newMetaInfo string,
) error {
	stmt := `
UPDATE assets
SET meta_info = $1
WHERE address = $2 AND layer_zero_chain_id = $3;`

	_, err := db.SQL.Exec(stmt, newMetaInfo, address, layerZeroChainID)
	if err != nil {
		return fmt.Errorf(
			"failed to update metadata for address %s on chain %d: %w",
			address, layerZeroChainID, err,
		)
	}
	return nil
}

// SaveStakerAsset saves a staker asset record in the database, including
// an entry in the history table.
func (db *Db) SaveStakerAsset(data *types.StakerAsset) error {
	if err := db.UpsertStakerAsset(data); err != nil {
		return err
	}
	if err := db.InsertStakerAssetHistory(data); err != nil {
		return err
	}
	return nil
}

// UpsertStakerAsset inserts or updates a staker asset record in the database.
func (db *Db) UpsertStakerAsset(data *types.StakerAsset) error {
	stmt := `
INSERT INTO staker_assets (staker_id, asset_id, deposited, free, delegated, pending_undelegation)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (staker_id, asset_id) DO UPDATE
SET deposited = EXCLUDED.deposited,
	free = EXCLUDED.free,
	delegated = EXCLUDED.delegated,
	pending_undelegation = EXCLUDED.pending_undelegation;`
	_, err := db.SQL.Exec(
		stmt, data.StakerID, data.AssetID,
		data.Deposited, data.Free, data.Delegated, data.PendingUndelegation,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert staker asset: %w", err)
	}
	return nil
}

// InsertStakerAssetHistory inserts a staker asset history record in the database.
func (db *Db) InsertStakerAssetHistory(data *types.StakerAsset) error {
	stmt := `
INSERT INTO staker_assets_history (staker_id, asset_id, deposited, free, delegated, pending_undelegation, block_height)
VALUES ($1, $2, $3, $4, $5, $6, $7);
    `
	_, err := db.SQL.Exec(
		stmt, data.StakerID, data.AssetID,
		data.Deposited, data.Free, data.Delegated, data.PendingUndelegation, data.Height,
	)
	if err != nil {
		return fmt.Errorf("failed to update staker asset history: %w", err)
	}
	return nil
}

// SaveOperatorAsset saves an operator asset record in the database, including
// an entry in the history table.
func (db *Db) SaveOperatorAsset(data *types.OperatorAsset) error {
	if err := db.UpsertOperatorAsset(data); err != nil {
		return err
	}
	if err := db.InsertOperatorAssetHistory(data); err != nil {
		return err
	}
	return nil
}

// UpsertOperatorAsset inserts or updates an operator asset record in the database.
func (db *Db) UpsertOperatorAsset(data *types.OperatorAsset) error {
	stmt := `
INSERT INTO operator_assets (operator, asset_id, delegated, pending_undelegation, share, self_share, delegated_share)
VALUES ($1, $2, $3, $4, $5, $6, $7);`
	_, err := db.SQL.Exec(
		stmt, data.OperatorAddress, data.AssetID,
		data.Delegated, data.PendingUndelegation,
		data.Share, data.SelfShare, data.DelegatedShare,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert operator asset: %w", err)
	}
	return nil
}

// InsertOperatorAssetHistory inserts an operator asset history record in the database.
func (db *Db) InsertOperatorAssetHistory(data *types.OperatorAsset) error {
	stmt := `
INSERT INTO operator_assets_history (operator, asset_id, delegated, pending_undelegation, share, self_share, delegated_share, block_height)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`
	_, err := db.SQL.Exec(
		stmt, data.OperatorAddress, data.AssetID,
		data.Delegated, data.PendingUndelegation,
		data.Share, data.SelfShare, data.DelegatedShare,
		data.Height,
	)
	if err != nil {
		return fmt.Errorf("failed to update operator asset history: %w", err)
	}
	return nil
}
