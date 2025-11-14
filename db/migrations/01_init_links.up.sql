-- Links table: holds canonical URL, short key (system or custom), lifecycle fields
CREATE TABLE IF NOT EXISTS links (
  id            BIGSERIAL PRIMARY KEY,
  key           VARCHAR(32) UNIQUE,                     -- nullable during create; set after base62(id)
  long_url      TEXT NOT NULL,                          -- canonicalized by app
  is_custom     BOOLEAN NOT NULL DEFAULT FALSE,         -- true for custom alias
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at    TIMESTAMPTZ NULL,
  is_disabled   BOOLEAN NOT NULL DEFAULT FALSE,
  CONSTRAINT key_format CHECK (key IS NULL OR key ~ '^[A-Za-z0-9_-]{3,32}$'),
  CONSTRAINT expires_after_created CHECK (expires_at IS NULL OR expires_at > created_at)
);

-- Idempotency for system-generated keys:
-- same canonical URL -> same system key (allow duplicates only for custom aliases)
CREATE UNIQUE INDEX IF NOT EXISTS uq_links_canonical_system
  ON links (long_url)
  WHERE is_custom = FALSE;

-- Optional: fast lookups for active keys (helps list/validation)
CREATE INDEX IF NOT EXISTS idx_links_key_active
  ON links (key)
  WHERE is_disabled = FALSE AND (expires_at IS NULL OR expires_at > NOW());