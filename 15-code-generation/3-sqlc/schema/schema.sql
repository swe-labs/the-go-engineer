-- Copyright (c) 2026 Rasel Hossen
-- Licensed under The Go Engineer License v1.0

-- schema.sql
CREATE TABLE users (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);
