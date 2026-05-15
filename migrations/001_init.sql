-- JobHub — complete schema (run this once on a fresh database)
-- For existing installations, run 002_update_existing.sql instead.

CREATE DATABASE IF NOT EXISTS job_portal CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE job_portal;

-- ─────────────────────────────────────────
-- USERS
-- ─────────────────────────────────────────
CREATE TABLE IF NOT EXISTS users (
    id                     INT AUTO_INCREMENT PRIMARY KEY,
    email                  VARCHAR(255) UNIQUE,
    mobile                 VARCHAR(20)  UNIQUE,
    password               VARCHAR(255),
    location               VARCHAR(100) NOT NULL,
    domain                 VARCHAR(100) NOT NULL,
    experience             VARCHAR(20)  DEFAULT '0-1',  -- e.g. "0-1", "1-3", "3-5", "5-10", "10+"
    notification_frequency ENUM('daily','weekly') DEFAULT 'daily',
    is_active              BOOLEAN      DEFAULT TRUE,
    is_verified            BOOLEAN      DEFAULT FALSE,
    verification_code      VARCHAR(12),
    created_at             TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    -- at least one of email or mobile must be provided
    CONSTRAINT chk_contact CHECK (email IS NOT NULL OR mobile IS NOT NULL)
);

-- ─────────────────────────────────────────
-- JOBS (sourced from LinkedIn, Indeed, Glassdoor via JSearch API)
-- ─────────────────────────────────────────
CREATE TABLE IF NOT EXISTS jobs (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    company     VARCHAR(255),
    location    VARCHAR(150),
    domain      VARCHAR(100),
    description TEXT,
    posted_at   TIMESTAMP,
    scraped_at  TIMESTAMP    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    source      VARCHAR(100),                      -- "LinkedIn", "Indeed", "Glassdoor", etc.
    source_url  VARCHAR(1000) UNIQUE,              -- UNIQUE prevents duplicate job entries
    INDEX idx_location  (location(50)),
    INDEX idx_domain    (domain),
    INDEX idx_posted_at (posted_at)
);

-- ─────────────────────────────────────────
-- SENT-JOB TRACKING  (prevents re-emailing the same job)
-- ─────────────────────────────────────────
CREATE TABLE IF NOT EXISTS user_job_sent (
    id      INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    job_id  INT NOT NULL,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_user_job (user_id, job_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (job_id)  REFERENCES jobs(id)  ON DELETE CASCADE
);
