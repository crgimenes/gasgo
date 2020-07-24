CREATE DATABASE gasgo
    WITH 
    OWNER = postgres
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1;

CREATE TABLE main
(
    id         uuid    NOT NULL,
    type_id    integer NOT NULL DEFAULT 0,
    value      text    NOT NULL,
    value_int  integer NOT NULL DEFAULT 0,
    parent_id  uuid,
    PRIMARY KEY (id)
);

