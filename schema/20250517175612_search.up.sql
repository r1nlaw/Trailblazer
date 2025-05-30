CREATE INDEX idx_name_fts ON landmark USING gin(to_tsvector('russian', name));

CREATE INDEX idx_address_fts ON landmark USING gin(to_tsvector('russian', address));
