-- init.sql
-- Create user if not exists
DO $$ BEGIN IF NOT EXISTS (
    SELECT
    FROM pg_user
    WHERE usename = 'postgres'
) THEN CREATE USER postgres;
END IF;
END $$;
ALTER USER postgres WITH SUPERUSER;
ALTER USER postgres WITH PASSWORD 'postgres';
CREATE DATABASE log_01;
CREATE DATABASE log_02;