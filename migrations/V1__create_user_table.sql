CREATE TABLE users
(
    id           uuid PRIMARY KEY,
    username     varchar(32),
    email        varchar(320),
    passwordHash varchar(256)
);