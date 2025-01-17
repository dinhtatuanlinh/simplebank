CREATE TABLE "sessions"
(
    "id"            uuid PRIMARY KEY,
    "username"      VARCHAR     NOT NULL,
    "refresh_token" VARCHAR     NOT NULL,
    "user_agent"    VARCHAR     NOT NULL,
    "client_ip"     VARCHAR     NOT NULL,
    "is_blocked"    BOOLEAN     NOT NULL DEFAULT false,
    "expires_at"    timestamptz NOT NULL,
    "created_at"    timestamptz NOT NULL DEFAULT (now())
);