-- Click events: exact analytics source of truth
CREATE TABLE IF NOT EXISTS clicks (
  id            BIGSERIAL PRIMARY KEY,
  link_id       BIGINT NOT NULL REFERENCES links(id) ON DELETE CASCADE,
  occurred_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  visitor_hash  TEXT NULL,                 -- hashed(IP + UA) or your chosen fingerprint
  country_code  CHAR(2) NULL,              -- optional; 'US', 'DE', etc. Use NULL or 'ZZ' if unknown
  user_agent    TEXT NULL,
  CONSTRAINT country_code_format CHECK (country_code IS NULL OR country_code ~ '^[A-Z]{2}$')
);

-- Time-ordered per-link access
CREATE INDEX IF NOT EXISTS idx_clicks_link_time
  ON clicks (link_id, occurred_at DESC);