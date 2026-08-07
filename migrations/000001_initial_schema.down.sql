DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS user_profiles;
DROP TABLE IF EXISTS users;

DO $$
BEGIN
    BEGIN
        DROP EXTENSION IF EXISTS pgcrypto;
    EXCEPTION
        WHEN dependent_objects_still_exist THEN
            NULL;
    END;
END;
$$;
