-- TODO: instead of assets_params, consider a mapping from chain id to gateway address
CREATE TABLE assets_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    -- int64 ctx.BlockHeight() fits in BIGINT
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

CREATE TABLE client_chains (
    name TEXT NOT NULL,
    meta_info TEXT NOT NULL,
    chain_id BIGINT,
    exocore_chain_index BIGINT,
    finalization_blocks BIGINT,
    layer_zero_chain_id BIGINT PRIMARY KEY,
    signature_type TEXT,
    address_length INT NOT NULL CHECK (address_length > 0)
);

-- rename from tokens to assets_tokens because tokens is used by pricefeed module
CREATE TABLE assets_tokens (
    -- generated for ease; not required to be part of the schema
    asset_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    symbol TEXT NOT NULL,
    address TEXT NOT NULL CHECK (address = lower(address)),
    decimals INT NOT NULL,
    layer_zero_chain_id BIGINT NOT NULL,
    exocore_chain_index BIGINT NOT NULL,
    meta_info TEXT,
    staking_total_amount NUMERIC NOT NULL DEFAULT 0,
    -- relational constraint
    CONSTRAINT fk_layer_zero_chain_id FOREIGN KEY (layer_zero_chain_id) REFERENCES client_chains (layer_zero_chain_id)
);
-- index by chain id
CREATE INDEX idx_tokens_layer_zero_chain_id ON assets_tokens (layer_zero_chain_id);

-- this is the latest state of a staker + asset id asset combination
-- it is kept to speed up queries as opposed to getting the latest from the history table
CREATE TABLE staker_assets (
    staker_id TEXT NOT NULL,
    asset_id TEXT NOT NULL,
    deposited NUMERIC NOT NULL DEFAULT 0,
    withdrawable NUMERIC NOT NULL DEFAULT 0,
    pending_undelegation NUMERIC NOT NULL DEFAULT 0,
    -- derived value via subtraction; only kept for speed
    delegated NUMERIC NOT NULL DEFAULT 0,
    -- not captured directly but via events
    lifetime_slashed NUMERIC NOT NULL DEFAULT 0,
    PRIMARY KEY (staker_id, asset_id),
    CONSTRAINT chk_total CHECK (deposited = withdrawable + pending_undelegation + delegated + lifetime_slashed),
    CONSTRAINT fk_asset_id FOREIGN KEY (asset_id) REFERENCES assets_tokens (asset_id)
    -- ideally, we would check that:
    -- (1) sum(total_deposited) <= tokens.staking_total_amount
    -- (2) staker_id.split("_")[1] == asset_id.split("_")[1]
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
    withdrawable NUMERIC,
    pending_undelegation NUMERIC,
    delegated NUMERIC,
    lifetime_slashed NUMERIC
) AS $$
BEGIN
    RETURN QUERY SELECT deposited, withdrawable, pending_undelegation, delegated, lifetime_slashed
    FROM staker_assets
    WHERE staker_id = p_staker_id AND asset_id = p_asset_id;
END;
$$ LANGUAGE plpgsql;

-- this table shared with x/operator and x/assets because
-- assets for an operator can only be tracked after the operator is created
CREATE TABLE operator_assets (
    operator_addr TEXT NOT NULL,
    asset_id TEXT NOT NULL,
    total_amount NUMERIC NOT NULL,
    pending_undelegation_amount NUMERIC NOT NULL,
    total_share NUMERIC NOT NULL,
    self_share NUMERIC NOT NULL DEFAULT 0,
    -- calculated / derived value
    other_share NUMERIC NOT NULL DEFAULT 0,
    PRIMARY KEY (operator_addr, asset_id),
    CONSTRAINT fk_asset_id FOREIGN KEY (asset_id) REFERENCES assets_tokens (asset_id),
    CONSTRAINT chk_total_share CHECK (total_share = self_share + other_share)
);
CREATE INDEX idx_operator_assets_operator ON operator_assets (operator_addr);
CREATE INDEX idx_operator_assets_asset_id ON operator_assets (asset_id);

CREATE OR REPLACE FUNCTION get_operator_assets(
    p_operator_addr TEXT,
    p_asset_id TEXT
) RETURNS TABLE (
    total_amount NUMERIC,
    pending_undelegation_amount NUMERIC,
    total_share NUMERIC,
    self_share NUMERIC,
    other_share NUMERIC
) AS $$
BEGIN
    RETURN QUERY SELECT total_amount, pending_undelegation_amount, total_share, self_share, other_share
    FROM operator_assets
    WHERE operator_addr = p_operator_addr AND asset_id = p_asset_id;
END;
$$ LANGUAGE plpgsql;

-- TODO add indexes by quantity?
