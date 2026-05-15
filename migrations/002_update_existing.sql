-- Run this if you already have a job_portal database from an earlier version.
-- Safe to run multiple times (uses IF NOT EXISTS / IGNORE).

USE job_portal;

-- Add columns that were missing from the original schema
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS mobile             VARCHAR(20),
    ADD COLUMN IF NOT EXISTS password           VARCHAR(255),
    ADD COLUMN IF NOT EXISTS experience         VARCHAR(20) DEFAULT '0-1',
    ADD COLUMN IF NOT EXISTS is_verified        BOOLEAN     DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS verification_code  VARCHAR(12);

-- Add source column and make source_url unique on jobs table
ALTER TABLE jobs
    ADD COLUMN IF NOT EXISTS source VARCHAR(100),
    MODIFY COLUMN source_url VARCHAR(1000);

-- Add unique constraint on source_url (skip if already exists)
SET @x = (SELECT COUNT(*) FROM information_schema.statistics
           WHERE table_schema = DATABASE() AND table_name = 'jobs' AND index_name = 'uq_source_url');
SET @sql = IF(@x = 0, 'ALTER TABLE jobs ADD UNIQUE INDEX uq_source_url (source_url(500))', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Remove old email_contact column if present (replaced by source_url)
SET @y = (SELECT COUNT(*) FROM information_schema.columns
           WHERE table_schema = DATABASE() AND table_name = 'jobs' AND column_name = 'email_contact');
SET @sql2 = IF(@y > 0, 'ALTER TABLE jobs DROP COLUMN email_contact', 'SELECT 1');
PREPARE stmt2 FROM @sql2; EXECUTE stmt2; DEALLOCATE PREPARE stmt2;
