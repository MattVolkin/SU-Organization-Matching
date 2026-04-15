-- 1. Extensions (Optional but recommended for IDs)
-- Using UUIDs is often better for security, but we'll stick to BIGINT for simplicity.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. User Table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    google_id VARCHAR(255) UNIQUE NOT NULL, -- Login ID from Google
    email VARCHAR(255) UNIQUE NOT NULL,
    tags JSONB DEFAULT '[]'::jsonb,
    genders JSONB NOT NULL DEFAULT '[]'::jsonb,
    ethnicities JSONB NOT NULL DEFAULT '[]'::jsonb,
    religions JSONB NOT NULL DEFAULT '[]'::jsonb,
    dedicated_majors JSONB NOT NULL DEFAULT '[]'::jsonb,
    other JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Club Table
CREATE TABLE IF NOT EXISTS clubs (
    id SERIAL PRIMARY KEY,
    club_name VARCHAR(255) NOT NULL,
    description TEXT,
    meeting_time VARCHAR(255),
    image_path VARCHAR(512), -- Store the path, not the file
    external_link VARCHAR(512),
    contact_info TEXT, -- Flexible for GroupMe, advisor emails, etc.
    include_officer_emails BOOLEAN DEFAULT FALSE,
    personality JSONB NOT NULL DEFAULT '[]'::jsonb,
    activities JSONB NOT NULL DEFAULT '[]'::jsonb,
    activities_description TEXT DEFAULT '',
    genders JSONB NOT NULL DEFAULT '[]'::jsonb,
    ethnicities JSONB NOT NULL DEFAULT '[]'::jsonb,
    religions JSONB NOT NULL DEFAULT '[]'::jsonb,
    strict_genders BOOLEAN NOT NULL DEFAULT FALSE,
    dedicated_majors JSONB NOT NULL DEFAULT '[]'::jsonb,
    associated_majors JSONB NOT NULL DEFAULT '[]'::jsonb,
    other JSONB NOT NULL DEFAULT '[]'::jsonb,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Keep existing databases in sync when re-running this script.
ALTER TABLE users DROP COLUMN IF EXISTS display_name;
ALTER TABLE users DROP COLUMN IF EXISTS personality;
ALTER TABLE users DROP COLUMN IF EXISTS activities;
ALTER TABLE users ADD COLUMN IF NOT EXISTS genders JSONB NOT NULL DEFAULT '[]'::jsonb;
ALTER TABLE users ADD COLUMN IF NOT EXISTS ethnicities JSONB NOT NULL DEFAULT '[]'::jsonb;
ALTER TABLE users ADD COLUMN IF NOT EXISTS religions JSONB NOT NULL DEFAULT '[]'::jsonb;
ALTER TABLE users ADD COLUMN IF NOT EXISTS dedicated_majors JSONB NOT NULL DEFAULT '[]'::jsonb;
ALTER TABLE users DROP COLUMN IF EXISTS associated_majors;
ALTER TABLE users ADD COLUMN IF NOT EXISTS other JSONB NOT NULL DEFAULT '[]'::jsonb;
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS meeting_time VARCHAR(255);
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS personality JSONB NOT NULL DEFAULT '[]'::jsonb;
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS activities JSONB NOT NULL DEFAULT '[]'::jsonb;
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS activities_description TEXT DEFAULT '';
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS genders JSONB NOT NULL DEFAULT '[]'::jsonb;
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS ethnicities JSONB NOT NULL DEFAULT '[]'::jsonb;
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS religions JSONB NOT NULL DEFAULT '[]'::jsonb;
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS strict_genders BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS dedicated_majors JSONB NOT NULL DEFAULT '[]'::jsonb;
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS associated_majors JSONB NOT NULL DEFAULT '[]'::jsonb;
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS other JSONB NOT NULL DEFAULT '[]'::jsonb;

-- 4. Club Leaders (Join Table for Many-to-Many)
CREATE TABLE IF NOT EXISTS club_leaders (
    club_id INTEGER REFERENCES clubs(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (club_id, user_id) -- Composite Primary Key (Implicit Index)
);

-- 5. Questions Table
CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    question_type VARCHAR(50) NOT NULL, -- e.g., 'multiple_choice', 'text'
    -- Optimized: Use JSONB for multi-language support (English, Spanish, etc.)
    -- Example structure: {"en": ["Creative", "A trait related to imagination and originality"]}
    translations JSONB NOT NULL, 
    is_active BOOLEAN DEFAULT TRUE
);

-- 6. Answers Table
CREATE TABLE IF NOT EXISTS answers (
    id SERIAL PRIMARY KEY,
    question_id INTEGER REFERENCES questions(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    answer_text TEXT NOT NULL,
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

---
--- PERFORMANCE & SECURITY OPTIMIZATIONS
---

-- Composite Index for faster Leader Lookups (if searching by user)
CREATE INDEX IF NOT EXISTS idx_club_leaders_user_id ON club_leaders(user_id);

-- GIN Index for the Translations JSONB column
-- This makes searching for specific language keys extremely fast.
CREATE INDEX IF NOT EXISTS idx_questions_translations ON questions USING GIN (translations);

-- Index for User Answers (Speeds up loading a user's specific profile)
CREATE INDEX IF NOT EXISTS idx_answers_user_id ON answers(user_id);

-- Permission Hardening
-- Ensure the 'dev_user' owns these tables
ALTER TABLE users OWNER TO dev_user;
ALTER TABLE clubs OWNER TO dev_user;
ALTER TABLE club_leaders OWNER TO dev_user;
ALTER TABLE questions OWNER TO dev_user;
ALTER TABLE answers OWNER TO dev_user;