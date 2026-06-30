-- Mentor profile: avatar and bio for personal center

ALTER TABLE mentors
    ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(512) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS bio TEXT NOT NULL DEFAULT '';
