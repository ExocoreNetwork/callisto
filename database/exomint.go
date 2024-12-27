package database

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/callisto/v4/types"
)

// SaveExomintParams allows to store the given params inside the database
func (db *Db) SaveExomintParams(params *types.ExomintParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling exomint params: %s", err)
	}

	stmt := `
INSERT INTO exomint_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE exomint_params.height <= excluded.height`

	_, err = db.SQL.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing exomint params: %s", err)
	}

	return nil
}

// AppendMintHistory allows to store the given mint history inside the database.
// The mint history can be updated for the same epoch, if there is a reorg, or
// there is a chain restart without retaining the epoch numbers (which we should
// not do). For such cases, the mint history will be updated. Otherwise, there
// should be no conflicts.
func (db *Db) AppendMintHistory(history *types.MintHistory) error {
	stmt := `
INSERT INTO exomint_history (block_height, quantity_minted, epoch_id, epoch_number, denom)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (epoch_id, epoch_number) DO UPDATE SET 
	block_height = EXCLUDED.block_height,
	quantity_minted = EXCLUDED.quantity_minted,
	denom = EXCLUDED.denom`
	_, err := db.SQL.Exec(stmt, history.Height, history.Amount, history.EpochID, history.EpochNumber, history.Denom)
	if err != nil {
		return fmt.Errorf("error while appending mint history: %s", err)
	}

	return nil
}
