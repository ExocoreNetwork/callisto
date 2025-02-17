-- TODO fully track x/avs level details
CREATE TABLE avs (
    avs_addr TEXT NOT NULL PRIMARY KEY
);

CREATE TABLE chain_id_to_avs_addr (
    chain_id TEXT NOT NULL PRIMARY KEY,
    avs_addr TEXT NOT NULL REFERENCES avs (avs_addr)
);

-- TODO vote power snapshots?
