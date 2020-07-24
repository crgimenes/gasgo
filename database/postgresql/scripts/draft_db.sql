CREATE TABLE main
(
    id             uuid    NOT NULL,
    metadata_id    integer NOT NULL DEFAULT 0,
    value          text    NOT NULL,
    value_int      integer NOT NULL DEFAULT 0,
    parent_id      uuid,
    PRIMARY KEY (id)
);


CREATE TABLE metadata
(
    id             integer NOT NULL,
    label          text    NOT NULL,
    PRIMARY KEY (id)
);
