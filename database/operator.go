package database

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/callisto/v4/types"
)

// SaveOperatorDetail saves the operator details into the database
func (db *Db) SaveOperatorDetail(operator *types.Operator) error {
	// convert time to the correct format
	parsed, err := sdk.ParseTime(operator.CommissionUpdateTime)
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}
	// Insert into operators table
	operatorStmt := `
INSERT INTO operators (earnings_addr, approve_addr, operator_meta_info, commission_rate, max_commission_rate, max_change_rate, commission_last_updated)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (earnings_addr) DO UPDATE
SET approve_addr = EXCLUDED.approve_addr,
	operator_meta_info = EXCLUDED.operator_meta_info,
	commission_rate = EXCLUDED.commission_rate,
	max_commission_rate = EXCLUDED.max_commission_rate,
	max_change_rate = EXCLUDED.max_change_rate,
	commission_last_updated = EXCLUDED.commission_last_updated;`
	_, err = db.SQL.Exec(
		operatorStmt,
		operator.EarningsAddress,
		operator.ApproveAddress,
		operator.MetaInfo,
		operator.Rate,
		operator.MaxRate,
		operator.MaxChangeRate,
		// use the parsed time to convert 1:1 to database TIMESTAMP format
		parsed,
	)
	if err != nil {
		return fmt.Errorf("failed to save operator details: %w", err)
	}

	// TODO: handle client chain earnings address list
	return nil
}

// SaveOperatorConsKey saves the operator consensus key into the database
func (db *Db) SaveOperatorConsKey(operatorAddr, chainID, pubkeyHex, consAddress string) error {
	stmt := `
INSERT INTO consensus_keys (operator_addr, chain_id, pubkey_hex, consensus_address)
VALUES ($1, $2, $3, $4)
ON CONFLICT (operator_addr, chain_id) DO UPDATE
SET consensus_address = EXCLUDED.consensus_address,
	pubkey_hex = EXCLUDED.pubkey_hex;`
	_, err := db.SQL.Exec(stmt, operatorAddr, chainID, pubkeyHex, consAddress)
	if err != nil {
		return fmt.Errorf("failed to save operator consensus key: %w", err)
	}
	return nil
}

// SaveOptedState inserts or updates an opted-in state into the avs_opt_ins table.
func (db *Db) SaveOptedState(data *types.Opted) error {
	// Prepare the SQL statement
	stmt := `
INSERT INTO avs_opt_ins (operator_addr, avs_addr, slash_contract, opt_in_height, opt_out_height, jailed)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (operator_addr, avs_addr) DO UPDATE
SET slash_contract = EXCLUDED.slash_contract,
	opt_in_height = EXCLUDED.opt_in_height,
	opt_out_height = EXCLUDED.opt_out_height,
	jailed = EXCLUDED.jailed;`

	// Execute the SQL statement
	_, err := db.SQL.Exec(
		stmt,
		data.OperatorAddress,
		data.AvsAddress,
		data.SlashContract,
		data.InHeight,
		data.OutHeight,
		data.Jailed,
	)

	if err != nil {
		return fmt.Errorf("failed to save opted state: %w", err)
	}

	return nil
}

// SaveOperatorUSDValue inserts or updates an operator USD value into the operator_usd_values table.
func (db *Db) SaveOperatorUSDValue(data *types.OperatorUSDValue) error {
	// Prepare the SQL statement
	stmt := `
INSERT INTO operator_usd_values (operator_addr, avs_addr, self_usd_value, total_usd_value, other_usd_value, active_usd_value)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (operator_addr, avs_addr) DO UPDATE
SET self_usd_value = EXCLUDED.self_usd_value,
	total_usd_value = EXCLUDED.total_usd_value,
	active_usd_value = EXCLUDED.active_usd_value;`

	// Execute the SQL statement
	_, err := db.SQL.Exec(
		stmt,
		data.OperatorAddress,
		data.AvsAddress,
		data.SelfUSDValue,
		data.TotalUSDValue,
		data.OtherUSDValue,
		data.ActiveUSDValue,
	)

	if err != nil {
		return fmt.Errorf("failed to save operator USD value: %w", err)
	}

	return nil
}

// DeleteOperatorUSDValue deletes an operator USD value from the operator_usd_values table.
func (db *Db) DeleteOperatorUSDValue(operatorAddr, avsAddr string) error {
	stmt := `
DELETE FROM operator_usd_values
WHERE operator_addr = $1 AND avs_addr = $2;`
	_, err := db.SQL.Exec(stmt, operatorAddr, avsAddr)
	if err != nil {
		return fmt.Errorf("failed to delete operator USD value: %w", err)
	}
	return nil
}

// SaveAvsUSDValue inserts or updates an AVS USD value into the avs_usd_values table.
func (db *Db) SaveAvsUSDValue(data *types.AvsUSDValue) error {
	// Prepare the SQL statement
	stmt := `
INSERT INTO avs_usd_values (avs_addr, usd_value)
VALUES ($1, $2)
ON CONFLICT (avs_addr) DO UPDATE
SET usd_value = EXCLUDED.usd_value;`

	// Execute the SQL statement
	_, err := db.SQL.Exec(
		stmt,
		data.AvsAddress,
		data.USDValue,
	)

	if err != nil {
		return fmt.Errorf("failed to save AVS USD value: %w", err)
	}

	return nil
}

// DeleteAvsUSDValue deletes an AVS USD value from the avs_usd_values table.
func (db *Db) DeleteAvsUSDValue(avsAddr string) error {
	stmt := `
DELETE FROM avs_usd_values
WHERE avs_addr = $1;`
	_, err := db.SQL.Exec(stmt, avsAddr)
	if err != nil {
		return fmt.Errorf("failed to delete AVS USD value: %w", err)
	}
	return nil
}

// SaveOperatorPrevConsKey inserts or updates an operator's previous consensus key into the consensus_keys table.
// TODO: it is not clear if this is even worth tracking.
func (db *Db) SaveOperatorPrevConsKey(operatorAddr, chainID, pubkeyHex, consAddress string) error {
	// Prepare the SQL statement to update previous consensus key fields
	stmt := `
UPDATE consensus_keys
SET prev_pubkey_hex = $3,
	prev_cons_addr = $4
WHERE operator_addr = $1 AND chain_id = $2;`

	// Execute the SQL statement
	_, err := db.SQL.Exec(
		stmt,
		operatorAddr,
		chainID,
		pubkeyHex,
		consAddress,
	)

	if err != nil {
		return fmt.Errorf("failed to save previous consensus key: %w", err)
	}

	return nil
}

// ClearOperatorPrevConsKey clears an operator's previous consensus key from the consensus_keys table.
func (db *Db) ClearOperatorPrevConsKey(operatorAddr, chainID string) error {
	stmt := `
UPDATE consensus_keys
SET prev_pubkey_hex = NULL,
	prev_cons_addr = NULL
WHERE operator_addr = $1 AND chain_id = $2;`
	_, err := db.SQL.Exec(stmt, operatorAddr, chainID)
	if err != nil {
		return fmt.Errorf("failed to clear operator previous consensus key: %w", err)
	}
	return nil
}

// MarkOperatorKeyRemoval marks an operator's key for removal by setting the is_removing field to true.
func (db *Db) MarkOperatorKeyRemoval(operatorAddr, chainID string) error {
	// Prepare the SQL statement to mark operator key removal
	stmt := `
UPDATE consensus_keys
SET is_removing = TRUE
WHERE operator_addr = $1 AND chain_id = $2;`

	// Execute the SQL statement
	_, err := db.SQL.Exec(
		stmt,
		operatorAddr,
		chainID,
	)

	if err != nil {
		return fmt.Errorf("failed to mark operator key removal: %w", err)
	}

	return nil
}

// RemoveOperatorConsKey removes an operator's consensus key from the consensus_keys table.
func (db *Db) RemoveOperatorConsKey(operatorAddr, chainID string) error {
	// Prepare the SQL statement to remove operator key
	stmt := `
DELETE FROM consensus_keys
WHERE operator_addr = $1 AND chain_id = $2;`

	// Execute the SQL statement
	_, err := db.SQL.Exec(
		stmt,
		operatorAddr,
		chainID,
	)

	if err != nil {
		return fmt.Errorf("failed to remove operator consensus key: %w", err)
	}

	return nil
}
