CREATE TABLE vehicle_locations (
    id         SERIAL PRIMARY KEY,
    vehicle_id TEXT,
    latitude   NUMERIC(10, 7),
    longitude  NUMERIC(10, 7),
    timestamp  BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM now())::BIGINT)
)