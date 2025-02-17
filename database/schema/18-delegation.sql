-- Remember that staker_assets and operator_assets are tracked in 15-assets.sql
-- We seperately track the impact of undelegations, native token delegations 
-- (because no staker asset) and consequently some slashing.

CREATE TABLE exo_asset_delegation (
    -- include the staker_id _0x0 suffix for ease of use with the other tables
    staker_id TEXT NOT NULL,
    operator_addr TEXT NOT NULL,
    -- derived from the delegation events
    delegated NUMERIC NOT NULL DEFAULT 0,
    -- derived from the undelegation events
    -- when an undelegation begins, add the amount to the pending_undelegation
    -- when it completes, subtract the ActualCompletedAmount from the pending_undelegation
    -- the difference between Amount and ActualCompletedAmount is the slashed amount
    -- but that will be emitted as a slashing event too
    pending_undelegation NUMERIC NOT NULL DEFAULT 0,
    -- cumulative, derived from the slashing events
    lifetime_slashed NUMERIC NOT NULL DEFAULT 0,
    -- no constraints on the amount because no deposits or withdrawals permitted
    PRIMARY KEY (staker_id, operator_addr),
    CONSTRAINT fk_operator FOREIGN KEY (operator_addr) REFERENCES operators (earnings_addr)
);

-- state<staker_id + asset_id + operator>
CREATE TABLE delegation_states (
    staker_id TEXT NOT NULL,
    asset_id TEXT NOT NULL,
    operator_addr TEXT NOT NULL,
    undelegatable_share NUMERIC NOT NULL DEFAULT 0,
    wait_undelegation_amount NUMERIC NOT NULL DEFAULT 0,
    PRIMARY KEY (staker_id, asset_id, operator_addr),
    CONSTRAINT fk_operator FOREIGN KEY (operator_addr) REFERENCES operators (earnings_addr),
    CONSTRAINT fk_asset_id FOREIGN KEY (asset_id) REFERENCES assets_tokens (asset_id)
);

-- staker to operator such that staker is unique
CREATE TABLE staker_operator_association (
    staker_id TEXT NOT NULL,
    operator_addr TEXT NOT NULL,
    -- each staker_id is associated with at most one operator
    PRIMARY KEY (staker_id),
    -- no contraint tying back to x/assets because it is possible to associate
    -- without any delegation. this is not subject to Sybil limitations because
    -- the downside is limited.
    CONSTRAINT fk_operator FOREIGN KEY (operator_addr) REFERENCES operators (earnings_addr)
);

-- no need for PRIMARY KEY index because it is automatically indexed
CREATE INDEX idx_association_by_operator ON staker_operator_association (operator_addr);

-- track the stakers for each operator and asset
CREATE TABLE operator_asset_stakers (
    operator_addr TEXT NOT NULL,
    asset_id TEXT NOT NULL,
    staker_id TEXT NOT NULL,
    PRIMARY KEY (operator_addr, asset_id, staker_id),
    CONSTRAINT fk_operator FOREIGN KEY (operator_addr) REFERENCES operators (earnings_addr),
    CONSTRAINT fk_asset FOREIGN KEY (asset_id) REFERENCES assets_tokens (asset_id)
);

CREATE INDEX idx_operator_asset_stakers_operator ON operator_asset_stakers (operator_addr);
CREATE INDEX idx_operator_asset_stakers_asset ON operator_asset_stakers (asset_id);
CREATE INDEX idx_operator_asset_stakers_operator_asset ON operator_asset_stakers (operator_addr, asset_id);
CREATE INDEX idx_operator_asset_stakers_staker ON operator_asset_stakers (staker_id);

CREATE TABLE undelegation_records (
    record_id TEXT PRIMARY KEY,
    staker_id TEXT NOT NULL,
    asset_id TEXT NOT NULL,
    operator_addr TEXT NOT NULL,
    tx_hash TEXT NOT NULL,
    block_number BIGINT NOT NULL,
    completed_epoch_identifier TEXT NOT NULL,
    completed_epoch_number BIGINT NOT NULL,
    undelegation_id BIGINT NOT NULL,
    amount NUMERIC NOT NULL,
    actual_completed_amount NUMERIC NOT NULL,
    hold_count BIGINT NOT NULL DEFAULT 0,
    -- the height at which it is matured, 0 if not matured
    maturity_height BIGINT NOT NULL DEFAULT 0,
    CONSTRAINT fk_operator FOREIGN KEY (operator_addr) REFERENCES operators (earnings_addr),
    CONSTRAINT fk_asset FOREIGN KEY (asset_id) REFERENCES assets_tokens (asset_id)
);

CREATE INDEX idx_undelegation_records_operator ON undelegation_records (operator_addr);
CREATE INDEX idx_undelegation_records_asset ON undelegation_records (asset_id);
CREATE INDEX idx_undelegation_records_staker ON undelegation_records (staker_id);
CREATE INDEX idx_undelegation_records_undelegation_id ON undelegation_records (undelegation_id);
