CREATE TYPE role AS ENUM('admin');

CREATE TABLE "user" (
    id uuid PRIMARY KEY,
    login VARCHAR(30) UNIQUE NOT NULL CHECK (length(login) > 0),
    name VARCHAR(30) NOT NULL CHECK (length(name) > 0),
    surname VARCHAR(30) NOT NULL CHECK (length(surname) > 0),
    phonenumber VARCHAR(12) NOT NULL CHECK (length(phonenumber) = 12),
    passhash bytea NOT NULL CHECK (length(passhash) = 60),
    email VARCHAR(60) NOT NULL CHECK (length(email) > 4),
    role role
);
