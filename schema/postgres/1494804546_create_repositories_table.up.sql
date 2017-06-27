CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE repositories (
  id             UUID PRIMARY KEY                           DEFAULT gen_random_uuid(),
  name           TEXT NOT NULL,
  description    TEXT,
  website        TEXT,
  default_branch TEXT         NOT NULL,
  private        BOOLEAN      NOT NULL                      DEFAULT TRUE,
  bare           BOOLEAN      NOT NULL                      DEFAULT TRUE,
  created_at     TIMESTAMPTZ  NOT NULL                      DEFAULT now(),
  updated_at     TIMESTAMPTZ  NOT NULL                      DEFAULT now(),
  owner_id       UUID REFERENCES users NOT NULL
);

ALTER TABLE repositories ADD UNIQUE (name, owner_id);
