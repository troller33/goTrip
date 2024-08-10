CREATE TABLE users (
    id       UUID        PRIMARY KEY,
    email    VARCHAR(33) UNIQUE NOT NULL,
    name     VARCHAR(33) NOT NULL,
    password VARCHAR(66) NOT NULL,
    admin    BOOLEAN     NOT NULL DEFAULT false
);


CREATE TABLE destination (
    id          UUID         PRIMARY KEY,
    name        VARCHAR(128) NOT NULL,
    description text         NOT NULL,
    attraction  text         NOT NULL
);

CREATE TABLE trip (
    id             UUID  PRIMARY KEY,
    name           text  NOT NULL,
    start_date     text  NOT NULL,
    end_date       text  NOT NULL,
    destination_id UUID  REFERENCES destination(id)
);
