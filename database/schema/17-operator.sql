CREATE TABLE operators (
    earnings_addr TEXT NOT NULL PRIMARY KEY,
    approve_addr TEXT NOT NULL,
    operator_meta_info TEXT,
    commission_rate NUMERIC NOT NULL,
    max_commission_rate NUMERIC NOT NULL,
    max_change_rate NUMERIC NOT NULL,
    -- we use ctx.BlockTime() which is in UTC, so drop the TZ
    commission_last_updated TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

-- I made this table a bit separate because it is not used actively yet.
CREATE TABLE client_chain_earning_addresses (
    operator_addr TEXT NOT NULL,
    lz_client_chain_id BIGINT NOT NULL,
    client_chain_earning_addr TEXT NOT NULL,
    PRIMARY KEY (operator_addr, lz_client_chain_id),
    CONSTRAINT fk_operator_addr FOREIGN KEY (operator_addr) REFERENCES operators (earnings_addr)
);

-- operator to avs mapping
CREATE TABLE operator_avs_opt_ins (
    operator_addr TEXT NOT NULL,
    avs_addr TEXT NOT NULL,
    -- can be null for x/dogfood, but not sure why it is tracked because it is an x/avs property
    slash_contract TEXT,
    -- not 0 because height is opted into
    opt_in_height BIGINT NOT NULL,
    -- default is max height
    opt_out_height BIGINT NOT NULL DEFAULT 18446744073709551615,
    jailed BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (operator_addr, avs_addr),
    CONSTRAINT fk_avs_addr FOREIGN KEY (avs_addr) REFERENCES avs (avs_addr),
    CONSTRAINT fk_operator_addr FOREIGN KEY (operator_addr) REFERENCES operators (earnings_addr)
);

-- this constraint can be added after the operators table is created
ALTER TABLE operator_assets
ADD CONSTRAINT fk_operator_addr
FOREIGN KEY (operator_addr) REFERENCES operators (earnings_addr);

CREATE TABLE consensus_keys (
    operator_addr TEXT NOT NULL,
    chain_id TEXT NOT NULL,
    pubkey_hex TEXT NOT NULL,
    cons_addr TEXT NOT NULL,
    -- optional, visible upon key rotation until pruned and thus can be NULL
    prev_pubkey_hex TEXT,
    prev_cons_addr TEXT,
    -- is_removing represents the situation in which the operator is in the process of opting out
    -- but has not yet completed the unbonding period
    is_removing BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (operator_addr, chain_id),
    CONSTRAINT fk_operator_addr FOREIGN KEY (operator_addr) REFERENCES operators (earnings_addr),
    CONSTRAINT fk_chain_id FOREIGN KEY (chain_id) REFERENCES chain_id_to_avs_addr (chain_id)
);

-- no correlation with NN-delegation.sql because the staker is not captured below.
CREATE TABLE operator_usd_values (
    operator_addr TEXT NOT NULL,
    avs_addr TEXT NOT NULL,
    self_usd_value NUMERIC NOT NULL,
    total_usd_value NUMERIC NOT NULL,
    -- for ease of query, we store the other usd value = total_usd_value - self_usd_value
    other_usd_value NUMERIC NOT NULL,
    -- 0 if self_usd_value < min_self_delegation for avs_addr
    active_usd_value NUMERIC NOT NULL,
    PRIMARY KEY (operator_addr, avs_addr),
    CONSTRAINT fk_operator_addr FOREIGN KEY (operator_addr) REFERENCES operators (earnings_addr),
    CONSTRAINT fk_avs_addr FOREIGN KEY (avs_addr) REFERENCES avs (avs_addr)
);

CREATE TABLE avs_usd_values (
    avs_addr TEXT NOT NULL PRIMARY KEY,
    usd_value NUMERIC NOT NULL,
    CONSTRAINT fk_avs_addr FOREIGN KEY (avs_addr) REFERENCES avs (avs_addr)
);

-- TODO slash states?
