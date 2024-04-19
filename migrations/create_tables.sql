BEGIN;
    -- DROP TABLE IF EXISTS Users;
    CREATE TABLE IF NOT EXISTS Users (
        id                  SERIAL 			PRIMARY KEY
        ,email              VARCHAR(255) 	not NULL UNIQUE
        ,firstName          varchar(255)	not NULL
        ,lastName        	  varchar(255)	not NULL
        ,password			      varchar(1000) 	not NULL
        ,datetimeCreated    TIMESTAMPTZ 	not NULL
        ,datetimeUpdated    TIMESTAMPTZ 	not NULL
    );
COMMIT;