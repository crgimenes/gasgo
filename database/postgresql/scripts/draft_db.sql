CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION UUID() RETURNS uuid AS $$
BEGIN
    RETURN uuid_generate_v4();
END;
$$ LANGUAGE 'plpgsql';

CREATE TABLE main
(
    id             uuid NOT NULL DEFAULT uuid(),
    parent_id      uuid,
    metadata_id    uuid NOT NULL,
    data           jsonb NOT NULL DEFAULT '{}',
    PRIMARY KEY (id)
);

CREATE TABLE metadata
(
    id             uuid NOT NULL DEFAULT uuid(),
    parent_id      uuid,
    type           character varying(255) NOT NULL,
    parameters     jsonb NOT NULL DEFAULT '{}',
    PRIMARY KEY (id)
);

ALTER TABLE main
    ADD CONSTRAINT metadata_fk FOREIGN KEY (metadata_id)
    REFERENCES metadata (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;