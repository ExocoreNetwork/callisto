CREATE TABLE assets_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

CREATE TABLE client_chains (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    meta_info TEXT NOT NULL,
    chain_id BIGINT NOT NULL,
    exocore_chain_index BIGINT NOT NULL,
    finalization_blocks BIGINT NOT NULL,
    layer_zero_chain_id BIGINT,
    signature_type TEXT,
    address_length INT NOT NULL
);

CREATE UNIQUE INDEX idx_client_chains_exocore_index ON client_chains (exocore_chain_index);

-- so far the primary identifier is the layer_zero_chain_id
CREATE UNIQUE INDEX idx_client_chains_layer_zero_chain_id ON client_chains (layer_zero_chain_id);

CREATE TABLE tokens (
    -- generated for ease; not required to be part of the schema
    asset_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    symbol TEXT NOT NULL,
    address TEXT NOT NULL,
    decimals INT NOT NULL,
    layer_zero_chain_id BIGINT NOT NULL,
    exocore_chain_index BIGINT NOT NULL,
    meta_info TEXT,
    -- staker asset state tracks post-slash deposits, but not pre-slash deposits
    -- those are tracked by the number below
    deposited NUMERIC NOT NULL DEFAULT 0,
    CONSTRAINT fk_layer_zero_chain_id FOREIGN KEY (layer_zero_chain_id) REFERENCES client_chains (layer_zero_chain_id),
    CONSTRAINT fk_exocore_chain_index FOREIGN KEY (exocore_chain_index) REFERENCES client_chains (exocore_chain_index)
);

CREATE INDEX idx_tokens_layer_zero_chain_id ON tokens (layer_zero_chain_id);

-- this is the latest state of a staker + asset id asset combination
-- it is kept to speed up queries as opposed to getting the latest from the history table
CREATE TABLE staker_assets (
    staker_id TEXT NOT NULL,
    asset_id TEXT NOT NULL,
    deposited NUMERIC NOT NULL DEFAULT 0,
    free NUMERIC NOT NULL DEFAULT 0,
    delegated NUMERIC NOT NULL DEFAULT 0,
    pending_undelegation NUMERIC NOT NULL DEFAULT 0,
    PRIMARY KEY (staker_id, asset_id),
    CONSTRAINT chk_total CHECK (deposited = free + delegated + pending_undelegation),
    CONSTRAINT fk_asset_id FOREIGN KEY (asset_id) REFERENCES tokens (asset_id)
    -- ideally, we would check that sum(total_deposited) <= tokens.staking_total_amount
    -- but that would require a trigger, which makes things more complicated than they need to be
    -- anyway, this is an indexer and not a business logic holder, so we can assume that the
    -- business logic is correct
);

CREATE INDEX idx_deposits_staker_id ON staker_assets (staker_id);
CREATE INDEX idx_deposits_asset_id ON staker_assets (asset_id);

CREATE OR REPLACE FUNCTION get_latest_staker_assets(
    p_staker_id TEXT,
    p_asset_id TEXT
) RETURNS TABLE (
    deposited NUMERIC,
    free NUMERIC,
    delegated NUMERIC,
    pending_undelegation NUMERIC
) AS $$
BEGIN
    RETURN QUERY SELECT deposited, free, delegated, pending_undelegation
    FROM staker_assets
    WHERE staker_id = p_staker_id AND asset_id = p_asset_id;
END;
$$ LANGUAGE plpgsql;

-- this is the history of a staker + asset id asset combination, indexed
-- by block height. it is the cumulative history, as a result of all
-- deposits and withdrawals. not just the change at the provided height.
-- it is trivial to calculate the change at a given height by using the
-- LAG function.
CREATE TABLE staker_assets_history (
    staker_id TEXT NOT NULL,
    asset_id TEXT NOT NULL,
    deposited NUMERIC NOT NULL DEFAULT 0,
    free NUMERIC NOT NULL DEFAULT 0,
    delegated NUMERIC NOT NULL DEFAULT 0,
    pending_undelegation NUMERIC NOT NULL DEFAULT 0,
    block_height BIGINT NOT NULL,
    PRIMARY KEY (staker_id, asset_id, block_height),
    CONSTRAINT chk_total CHECK (deposited = free + delegated + pending_undelegation),
    CONSTRAINT fk_asset_id FOREIGN KEY (asset_id) REFERENCES tokens (asset_id)
);

CREATE INDEX idx_staker_assets_history_staker_id ON staker_assets_history (staker_id);
CREATE INDEX idx_staker_assets_history_asset_id ON staker_assets_history (asset_id);
CREATE INDEX idx_staker_assets_history_block_height ON staker_assets_history (block_height);

CREATE OR REPLACE FUNCTION get_staker_assets_history(
    p_staker_id TEXT,
    p_asset_id TEXT
) RETURNS TABLE (
    block_height BIGINT,
    deposited NUMERIC,
    free NUMERIC,
    delegated NUMERIC,
    pending_undelegation NUMERIC
) AS $$
BEGIN
    RETURN QUERY SELECT block_height, deposited, free, delegated, pending_undelegation
    FROM staker_assets_history
    WHERE staker_id = p_staker_id AND asset_id = p_asset_id
    ORDER BY block_height ASC;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_staker_assets_history_at_height(
    p_staker_id TEXT,
    p_asset_id TEXT,
    p_block_height BIGINT
) RETURNS TABLE (
    block_height BIGINT,
    deposited NUMERIC,
    free NUMERIC,
    delegated NUMERIC,
    pending_undelegation NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT block_height, deposited, free, delegated, pending_undelegation
    FROM staker_assets_history
    WHERE staker_id = p_staker_id AND asset_id = p_asset_id AND block_height <= p_block_height
    ORDER BY block_height DESC
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE operator_assets (
    operator TEXT NOT NULL,
    asset_id TEXT NOT NULL,
    delegated NUMERIC NOT NULL, -- includes self-delegation, which is tracked via shares only
    pending_undelegation NUMERIC NOT NULL DEFAULT 0, -- includes self undelegation
    share NUMERIC NOT NULL,
    self_share NUMERIC NOT NULL,
    delegated_share NUMERIC NOT NULL DEFAULT 0,
    PRIMARY KEY (operator, asset_id),
    CONSTRAINT fk_asset_id FOREIGN KEY (asset_id) REFERENCES tokens (asset_id),
    CONSTRAINT chk_total_share CHECK (share = self_share + delegated_share)
);
CREATE INDEX idx_operator_assets_operator ON operator_assets (operator);
CREATE INDEX idx_operator_assets_asset_id ON operator_assets (asset_id);

CREATE OR REPLACE FUNCTION get_operator_assets(
    p_operator TEXT,
    p_asset_id TEXT
) RETURNS TABLE (
    delegated NUMERIC,
    pending_undelegation NUMERIC,
    share NUMERIC,
    self_share NUMERIC,
    delegated_share NUMERIC
) AS $$
BEGIN
    RETURN QUERY SELECT delegated, pending_undelegation, share, self_share, delegated_share
    FROM operator_assets
    WHERE operator = p_operator AND asset_id = p_asset_id;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE operator_assets_history (
    operator TEXT NOT NULL,
    asset_id TEXT NOT NULL,
    delegated NUMERIC NOT NULL,
    pending_undelegation NUMERIC NOT NULL DEFAULT 0,
    share NUMERIC NOT NULL,
    self_share NUMERIC NOT NULL,
    delegated_share NUMERIC NOT NULL DEFAULT 0,
    block_height BIGINT NOT NULL,
    PRIMARY KEY (operator, asset_id, block_height),
    CONSTRAINT fk_asset_id FOREIGN KEY (asset_id) REFERENCES tokens (asset_id),
    CONSTRAINT chk_total_share CHECK (share = self_share + delegated_share)
);
CREATE INDEX idx_operator_assets_history_operator ON operator_assets_history (operator);
CREATE INDEX idx_operator_assets_history_asset_id ON operator_assets_history (asset_id);
CREATE INDEX idx_operator_assets_history_block_height ON operator_assets_history (block_height);

CREATE OR REPLACE FUNCTION get_operator_assets_history(
    p_operator TEXT,
    p_asset_id TEXT
) RETURNS TABLE (
    block_height BIGINT,
    delegated NUMERIC,
    pending_undelegation NUMERIC,
    share NUMERIC,
    self_share NUMERIC,
    delegated_share NUMERIC
) AS $$
BEGIN
    RETURN QUERY SELECT block_height, delegated, pending_undelegation, share, self_share, delegated_share
    FROM operator_assets_history
    WHERE operator = p_operator AND asset_id = p_asset_id
    ORDER BY block_height ASC;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_operator_assets_history_at_height(
    p_operator TEXT,
    p_asset_id TEXT,
    p_block_height BIGINT
) RETURNS TABLE (
    delegated NUMERIC,
    pending_undelegation NUMERIC,
    share NUMERIC,
    self_share NUMERIC,
    delegated_share NUMERIC
) AS $$
BEGIN
    RETURN QUERY SELECT block_height, delegated, pending_undelegation, share, self_share, delegated_share
    FROM operator_assets_history
    WHERE operator = p_operator AND asset_id = p_asset_id AND block_height <= p_block_height
    ORDER BY block_height DESC
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;

-- TODO add indexes by quantity?