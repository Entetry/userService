CREATE TABLE users
(
    id           uuid PRIMARY KEY,
    username     varchar(32),
    email        varchar(320),
    passwordHash varchar(256),
    CONSTRAINT email_unique UNIQUE (email)
);

CREATE UNIQUE INDEX username_unique ON users(username);