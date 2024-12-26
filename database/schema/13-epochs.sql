-- static, unchanging data
CREATE TABLE epoch_definitions (
    identifier TEXT PRIMARY KEY,
    start_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    duration INTERVAL NOT NULL
);

-- dynamic data
CREATE TABLE epoch_states (
    identifier TEXT PRIMARY KEY REFERENCES epoch_definitions(identifier),
    current_epoch BIGINT NOT NULL,
    current_epoch_start_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    epoch_counting_started BOOLEAN NOT NULL DEFAULT FALSE,
    current_epoch_start_height BIGINT NOT NULL
);

-- These indexes help the functions below
CREATE INDEX idx_epoch_identifier_time ON epoch_states (identifier, current_epoch_start_time);
CREATE INDEX idx_epoch_identifier_height ON epoch_states (identifier, current_epoch_start_height);

-- This function is used to get the epoch number by height
CREATE OR REPLACE FUNCTION get_epoch_number_by_height(identifier TEXT, block_height BIGINT)
RETURNS BIGINT AS $$
    SELECT current_epoch
    FROM epoch_states
    WHERE identifier = $1
      AND current_epoch_start_height <= $2
    ORDER BY current_epoch_start_height DESC
    LIMIT 1;
$$ LANGUAGE sql STABLE;

-- This function is used to get the epoch number by time
CREATE OR REPLACE FUNCTION get_epoch_number_by_time(identifier TEXT, block_time TIMESTAMP)
RETURNS BIGINT AS $$
    SELECT current_epoch
    FROM epoch_states
    WHERE identifier = $1
      AND current_epoch_start_time <= $2
    ORDER BY current_epoch_start_time DESC
    LIMIT 1;
$$ LANGUAGE sql STABLE;
