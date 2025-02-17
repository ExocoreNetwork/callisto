package database

import "fmt"

// SaveAvsAddr saves the given avs address into the database.
// TODO: implement a full-fledged function that saves all AVS information
func (db *Db) SaveAvsAddr(avsAddr string) error {
	stmt := `
INSERT INTO avs (avs_addr)
VALUES ($1)
ON CONFLICT (avs_addr) DO NOTHING;`
	_, err := db.SQL.Exec(stmt, avsAddr)
	if err != nil {
		return fmt.Errorf("error while saving avs address: %s", err)
	}
	return nil
}

// SaveChainIdToAvsAddr saves the given chain id to avs address mapping into the database.
func (db *Db) SaveChainIdToAvsAddr(chainID, avsAddr string) error {
	stmt := `
INSERT INTO chain_id_to_avs_addr (chain_id, avs_addr)
VALUES ($1, $2)
ON CONFLICT (chain_id) DO UPDATE
SET avs_addr = EXCLUDED.avs_addr;`
	_, err := db.SQL.Exec(stmt, chainID, avsAddr)
	if err != nil {
		return fmt.Errorf("error while saving chain id to avs address: %s", err)
	}
	return nil
}
