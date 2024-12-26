package database

import (
	"fmt"

	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"
)

// SaveEpochs saves the given epochs into the database. This includes
// the definitions and the states of the epochs.
func (db *Db) SaveEpochs(epochs []epochstypes.EpochInfo) error {
	// first, we save the static data
	if err := db.SaveEpochDefinitions(epochs); err != nil {
		// already wrapped
		return err
	}
	// then we save the dynamic data
	return db.SaveEpochStates(epochs)
}

// SaveEpochDefinitions saves the given epochs definitions into the database.
// It is designed to be only called at genesis.
func (db *Db) SaveEpochDefinitions(epochs []epochstypes.EpochInfo) error {
	stmt := `
INSERT INTO epoch_definitions (identifier, start_time, duration) 
VALUES ($1, $2, $3)
ON CONFLICT (identifier) DO NOTHING;`
	for _, epoch := range epochs {
		_, err := db.SQL.Exec(stmt, epoch.Identifier, epoch.StartTime, epoch.Duration)
		if err != nil {
			return fmt.Errorf("error while saving epoch definitions: %s", err)
		}
	}
	return nil
}

// SaveEpochStates saves the given x/epochs states into the database.
func (db *Db) SaveEpochStates(epochs []epochstypes.EpochInfo) error {
	for _, epoch := range epochs {
		if err := db.SaveEpochState(epoch); err != nil {
			return err
		}
	}
	return nil
}

// SaveEpochState saves the given x/epochs state into the database.
func (db *Db) SaveEpochState(epoch epochstypes.EpochInfo) error {
	// the function should ideally be called only on epoch change.
	// however, we guard against that anyway.
	stmt := `
INSERT INTO epoch_states (identifier, current_epoch, current_epoch_start_time, epoch_counting_started, current_epoch_start_height) 
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (identifier) DO UPDATE 
SET 
    current_epoch = EXCLUDED.current_epoch,
    current_epoch_start_time = EXCLUDED.current_epoch_start_time,
    current_epoch_start_height = EXCLUDED.current_epoch_start_height,
    epoch_counting_started = EXCLUDED.epoch_counting_started
WHERE 
    EXCLUDED.current_epoch > epoch_states.current_epoch;`
	_, err := db.SQL.Exec(
		stmt, epoch.Identifier, epoch.CurrentEpoch, epoch.CurrentEpochStartTime,
		epoch.EpochCountingStarted, epoch.CurrentEpochStartHeight,
	)
	if err != nil {
		return fmt.Errorf("error while saving epoch state: %s", err)
	}
	return nil
}
