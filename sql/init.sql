CREATE TABLE IF NOT EXISTS users_data
(
    id          SERIAL PRIMARY KEY,
    user_id     INT         NOT NULL,
    service     VARCHAR(64) NOT NULL,
    login       VARCHAR(64) NOT NULL,
    password    VARCHAR(64) NOT NULL,
    last_update TIMESTAMP   NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ON users_data (user_id, service);
