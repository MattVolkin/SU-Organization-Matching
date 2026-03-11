-- 1. Extensions (Optional but recommended for IDs)
-- Using UUIDs is often better for security, but we'll stick to BIGINT for simplicity.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. User Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    google_id VARCHAR(255) UNIQUE NOT NULL, -- Login ID from Google
    email VARCHAR(255) UNIQUE NOT NULL,
    tags JSONB DEFAULT '[]'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Club Table
CREATE TABLE clubs (
    id SERIAL PRIMARY KEY,
    club_name VARCHAR(255) NOT NULL,
    description TEXT,
    meeting_time VARCHAR(255),
    image_path VARCHAR(512), -- Store the path, not the file
    external_link VARCHAR(512),
    contact_info TEXT, -- Flexible for GroupMe, advisor emails, etc.
    include_officer_emails BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Keep existing databases in sync when re-running this script.
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS meeting_time VARCHAR(255);

-- 4. Club Leaders (Join Table for Many-to-Many)
CREATE TABLE club_leaders (
    club_id INTEGER REFERENCES clubs(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (club_id, user_id) -- Composite Primary Key (Implicit Index)
);

-- 5. Questions Table
CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    question_type VARCHAR(50) NOT NULL, -- e.g., 'multiple_choice', 'text'
    -- Optimized: Use JSONB for multi-language support (English, Spanish, etc.)
    -- Example structure: {"en": "What is your name?", "es": "Cual es tu nombre?"}
    translations JSONB NOT NULL, 
    is_active BOOLEAN DEFAULT TRUE
);

-- 6. Answers Table
CREATE TABLE answers (
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
CREATE INDEX idx_club_leaders_user_id ON club_leaders(user_id);

-- GIN Index for the Translations JSONB column
-- This makes searching for specific language keys extremely fast.
CREATE INDEX idx_questions_translations ON questions USING GIN (translations);

-- Index for User Answers (Speeds up loading a user's specific profile)
CREATE INDEX idx_answers_user_id ON answers(user_id);

-- Permission Hardening
-- Ensure the 'dev_user' owns these tables
ALTER TABLE users OWNER TO dev_user;
ALTER TABLE clubs OWNER TO dev_user;
ALTER TABLE club_leaders OWNER TO dev_user;
ALTER TABLE questions OWNER TO dev_user;
ALTER TABLE answers OWNER TO dev_user;