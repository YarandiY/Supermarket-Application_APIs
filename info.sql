CREATE DATABASE my_database
  WITH OWNER = postgres
       ENCODING = 'UTF8'
       TABLESPACE = pg_default
       CONNECTION LIMIT = -1;

CREATE TABLE customers
(
 id serial NOT NULL,
 name character varying NOT NULL,
 tel bigint,
 address character varying,
 registerDate date,
 CONSTRAINT pk_customers PRIMARY KEY (id )
)
WITH (
 OIDS=FALSE
);
ALTER TABLE customers
 OWNER TO postgres;