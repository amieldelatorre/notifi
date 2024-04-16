-- SELECT * FROM pg_catalog.pg_tables;


BEGIN;
    CREATE TABLE IF NOT EXISTS Users (
        id                  SERIAL 			PRIMARY KEY
        ,email              VARCHAR(255) 	not NULL UNIQUE
        ,firstName          varchar(255)	not NULL
        ,lastName        	varchar(255)	not NULL
        ,password			varchar(1000) 	not NULL
        ,datetime_created   TIMESTAMPTZ 	not NULL
        ,datetime_updated   TIMESTAMPTZ 	not NULL
    );

-- ROLLBACK;
COMMIT;