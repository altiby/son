CREATE TABLE "users"
(
    "id"          TEXT PRIMARY KEY NOT NULL,
    "password"    TEXT             NOT NULL,
    "role"        TEXT,
    "first_name"  varchar(20),
    "second_name" varchar(20),
    "birthdate"   TEXT,
    "biography"   TEXT,
    "city"        TEXT,
    "created_at"  timestamp        NOT NULL DEFAULT 'now()'
);