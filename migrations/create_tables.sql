DO $$
BEGIN
    BEGIN
        CREATE TABLE IF NOT EXISTS Users (
            id                  INT             GENERATED ALWAYS AS IDENTITY PRIMARY KEY
            ,email              VARCHAR(255) 	NOT NULL UNIQUE
            ,firstName          VARCHAR(255)    NOT NULL
            ,lastName        	VARCHAR(255)    NOT NULL
            ,password			VARCHAR(1000)   NOT NULL
            ,datetimeCreated    TIMESTAMPTZ     NOT NULL
            ,datetimeUpdated    TIMESTAMPTZ     NOT NULL
        );

        CREATE TABLE IF NOT EXISTS Destinations (
            id                      INT             GENERATED ALWAYS AS IDENTITY PRIMARY KEY
            ,userId                 INT             NOT NULL
            ,type                   VARCHAR(255)    NOT NULL
            ,identifier             TEXT            NOT NULL
            ,datetimeCreated        TIMESTAMPTZ     NOT NULL
            ,datetimeUpdated        TIMESTAMPTZ     NOT NULL
            ,CONSTRAINT fk_userid
                FOREIGN KEY(userId)
                REFERENCES Users(id)
        );

        CREATE TABLE IF NOT EXISTS Messages (
            id                      INT             GENERATED ALWAYS AS IDENTITY PRIMARY KEY
            ,userId                 INT             NOT NULL
            ,destinationId          INT             NOT NULL
            ,title                  TEXT            NOT NULL
            ,body                   TEXT            NOT NULL
            ,status                 VARCHAR(255)    NOT NULL
            ,datetimeCreated        TIMESTAMPTZ     NOT NULL
            ,datetimeSendAttempt    TIMESTAMPTZ     NOT NULL
            ,CONSTRAINT fk_userid
                FOREIGN KEY(userId)
                REFERENCES Users(id)
            ,CONSTRAINT fk_destinationId
                FOREIGN KEY(destinationId)
                REFERENCES Destinations(id)
        );
        
    EXCEPTION
        WHEN OTHERS THEN
            RAISE NOTICE 'Error: %', SQLERRM;
            ROLLBACK; -- Rollback the entire transaction
    END;
    
    COMMIT;
END $$;