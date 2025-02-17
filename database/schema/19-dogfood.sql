CREATE TABLE dogfood_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    height     BIGINT  NOT NULL,
    epochs_until_unbonded BIGINT NOT NULL,
    epoch_identifier TEXT NOT NULL,
    max_validators BIGINT NOT NULL,
    historical_entries BIGINT NOT NULL,
    min_self_delegation NUMERIC NOT NULL,
    -- this is acceptable because it rarely changes. otherwise TEXT would be better, but then
    -- you'd need more rows in this table. so i think this is a good tradeoff.
    asset_ids TEXT[] NOT NULL,
    CHECK (one_row_id)
);

-- the validator table is created in 00-cosmos.sql; however, it does not contain the vote power
-- that table is fed via Tendermint responses, so it covers a list of validators with addrs
-- and pubkeys, but not their voting power. we create a new table to make up the difference.
-- it is simply a copy of the table created in 03-staking.sql
CREATE TABLE validator_voting_power
(
    validator_address TEXT   NOT NULL REFERENCES validator (consensus_address) PRIMARY KEY,
    voting_power      BIGINT NOT NULL,
    height            BIGINT NOT NULL REFERENCES block (height)
);
CREATE INDEX validator_voting_power_height_index ON validator_voting_power (height);

CREATE TABLE opt_out_expiries
(
    epoch_number BIGINT NOT NULL,
    operator_addr TEXT NOT NULL,
    completion_height BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (operator_addr),
    CONSTRAINT fk_operator_addr FOREIGN KEY (operator_addr) REFERENCES operators (earnings_addr)
    -- no constraint to refer back to the epoch identifier because
    -- it is only tracked in params and not here.
);

CREATE TABLE consensus_addrs_to_prune
(
    epoch_number BIGINT NOT NULL,
    consensus_addr TEXT NOT NULL,
    completion_height BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (consensus_addr),
    CONSTRAINT fk_consensus_addr FOREIGN KEY (consensus_addr) REFERENCES validator (consensus_address)
    -- no foreign key constraints here because removal of consensus_addr
);

CREATE TABLE undelegation_maturities
(
    epoch_number BIGINT NOT NULL,
    undelegation_record_id TEXT NOT NULL,
    completion_height BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (undelegation_record_id),
    CONSTRAINT fk_undelegation_record_id FOREIGN KEY (undelegation_record_id) REFERENCES undelegation_records (record_id)
);

CREATE TABLE last_total_power
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    total_power NUMERIC NOT NULL,
    CHECK (one_row_id)
);
