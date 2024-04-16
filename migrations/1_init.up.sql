CREATE TABLE IF NOT EXISTS users (
                                     id                BIGSERIAL  PRIMARY KEY,
                                     email             TEXT     NOT NULL UNIQUE,
                                     pass_hash BYTEA     NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE IF NOT EXISTS apps (
                                    id          BIGSERIAL PRIMARY KEY,
                                    name        TEXT NOT NULL UNIQUE,
                                    secret TEXT NOT NULL UNIQUE
);

-- регитрируем в таблице apps наше приложение с миграциями
INSERT INTO apps (name, secret) VALUES ('jwt_auth', 'secret_key');