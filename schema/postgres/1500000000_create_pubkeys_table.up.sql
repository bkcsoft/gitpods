CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE pubkeys (
  id          UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
  name        VARCHAR(32)  NOT NULL,
  fingerprint VARCHAR(100) NOT NULL,
  content     VARCHAR(500) NOT NULL,
  created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);
