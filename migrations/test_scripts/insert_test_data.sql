DO $$
BEGIN
    BEGIN
        INSERT INTO Users (email, firstName, lastName, password, datetimeCreated, datetimeUpdated) 
		VALUES ('isaac.newton@invalid.com', 'Isaac', 'Newton', '$argon2id$v=19$m=65536,t=1,p=11$V3dLd/hNnClN0U9mSu3IbQ$mcK9nxoJpUkWsWNbhX14tEC4pXp0oihqJQVePj7FFIc', NOW(), NOW());
        
        INSERT INTO Users (email, firstName, lastName, password, datetimeCreated, datetimeUpdated) 
		VALUES ('alberteinstein@example.invalid', 'Albert', 'Einstein', '$argon2id$v=19$m=65536,t=1,p=11$PQT2VdnSXGtAjLKmLHk7jA$hrclADmr/RTFGZgX0J2ujMmZg0adxhOOJczzp1YFMBk', NOW(), NOW());

        INSERT INTO Destinations (userId, type, identifier, datetimeCreated, datetimeUpdated) 
        VALUES (1, 'DISCORD', 'https://one.example.discord.webhook.invalid', NOW(), NOW());

        INSERT INTO Messages (userId, destinationId, title, body, status, datetimeCreated, datetimeSendAttempt)
        VALUES (1, 1, 'MessageTitle', 'MessageBody', 'PENDING', NOW(), NOW());
    EXCEPTION
        WHEN OTHERS THEN
            RAISE NOTICE 'Error: %', SQLERRM;
            ROLLBACK; -- Rollback the entire transaction
    END;
    
    COMMIT;
END $$;