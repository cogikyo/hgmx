package auth

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/cogikyo/hgmx/migrate"
)

// Track is the auth schema; list it before the site's own track.
var Track = migrate.Track{Name: "auth", Steps: []migrate.Step{
	{Version: 1, Name: "accounts", Up: accountsUp, Down: accountsDown},
	{Version: 2, Name: "usernames", Up: usernamesUp, Down: usernamesDown},
}}

func accountsUp(_ context.Context, db *gorm.DB) error {
	return db.Exec(`
CREATE TABLE auth_users (
	id uuid PRIMARY KEY,
	email text NOT NULL UNIQUE,
	handle bytea NOT NULL UNIQUE,
	created timestamptz NOT NULL
);

CREATE TABLE auth_passkeys (
	id bytea PRIMARY KEY,
	user_id uuid NOT NULL REFERENCES auth_users ON DELETE CASCADE,
	name text NOT NULL,
	credential jsonb NOT NULL,
	created timestamptz NOT NULL,
	used timestamptz
);
CREATE INDEX auth_passkeys_user ON auth_passkeys (user_id);

CREATE TABLE auth_codes (
	id uuid PRIMARY KEY,
	email text NOT NULL,
	purpose text NOT NULL CHECK (purpose IN ('signin', 'link')),
	mac bytea NOT NULL,
	attempts integer NOT NULL DEFAULT 0,
	created timestamptz NOT NULL,
	expires timestamptz NOT NULL,
	consumed timestamptz
);
CREATE INDEX auth_codes_email ON auth_codes (email, created);
CREATE INDEX auth_codes_created ON auth_codes (created);

CREATE TABLE auth_ceremonies (
	id bytea PRIMARY KEY,
	kind text NOT NULL CHECK (kind IN ('login', 'register')),
	session jsonb NOT NULL,
	user_id uuid REFERENCES auth_users ON DELETE CASCADE,
	library text NOT NULL,
	expires timestamptz NOT NULL
);

CREATE TABLE auth_sessions (
	id uuid PRIMARY KEY,
	user_id uuid NOT NULL REFERENCES auth_users ON DELETE CASCADE,
	token_hash bytea NOT NULL UNIQUE,
	created timestamptz NOT NULL,
	seen timestamptz NOT NULL,
	expires timestamptz NOT NULL
);
CREATE INDEX auth_sessions_user ON auth_sessions (user_id);
`).Error
}

func accountsDown(_ context.Context, db *gorm.DB) error {
	return db.Exec(`
DROP TABLE auth_sessions;
DROP TABLE auth_ceremonies;
DROP TABLE auth_codes;
DROP TABLE auth_passkeys;
DROP TABLE auth_users;
`).Error
}

func usernamesUp(_ context.Context, db *gorm.DB) error {
	return db.Exec(`
ALTER TABLE auth_users ALTER COLUMN email DROP NOT NULL, ADD COLUMN username text;
CREATE UNIQUE INDEX auth_users_username ON auth_users (lower(username COLLATE "C"));

ALTER TABLE auth_codes
	ADD COLUMN session_id uuid,
	DROP CONSTRAINT auth_codes_purpose_check,
	ADD CONSTRAINT auth_codes_purpose_check CHECK (purpose IN ('signin', 'link', 'verify'));
`).Error
}

func usernamesDown(_ context.Context, db *gorm.DB) error {
	return db.Exec(`
DELETE FROM auth_codes WHERE purpose = 'verify';
ALTER TABLE auth_codes
	DROP COLUMN session_id,
	DROP CONSTRAINT auth_codes_purpose_check,
	ADD CONSTRAINT auth_codes_purpose_check CHECK (purpose IN ('signin', 'link'));
ALTER TABLE auth_users DROP COLUMN username, ALTER COLUMN email SET NOT NULL;
`).Error
}

// sweep runs outside address locks; SKIP LOCKED avoids waiting on active rows.
func sweep(db *gorm.DB, library string) error {
	now := time.Now()
	for _, q := range []struct {
		sql  string
		args []any
	}{
		// Keep even consumed codes for 25 hours of per-address rate accounting.
		{"DELETE FROM auth_codes WHERE id IN (SELECT id FROM auth_codes WHERE created < ? FOR UPDATE SKIP LOCKED)", []any{now.Add(-retention)}},
		{"DELETE FROM auth_sessions WHERE id IN (SELECT id FROM auth_sessions WHERE expires <= ? OR seen <= ? FOR UPDATE SKIP LOCKED)", []any{now, now.Add(-idleLife)}},
		{"DELETE FROM auth_ceremonies WHERE id IN (SELECT id FROM auth_ceremonies WHERE expires <= ? OR library <> ? FOR UPDATE SKIP LOCKED)", []any{now, library}},
	} {
		if err := db.Exec(q.sql, q.args...).Error; err != nil {
			return err
		}
	}
	return nil
}
