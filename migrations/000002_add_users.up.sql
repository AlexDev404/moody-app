-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    hashed_password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Add user_id to mood_entries table
ALTER TABLE mood_entries ADD COLUMN IF NOT EXISTS user_id UUID REFERENCES users(id);

-- Add user_id to playlists table
ALTER TABLE playlists ADD COLUMN IF NOT EXISTS user_id UUID REFERENCES users(id);

-- Create indexes for improved performance
CREATE INDEX IF NOT EXISTS idx_mood_entries_user_id ON mood_entries(user_id);
CREATE INDEX IF NOT EXISTS idx_playlists_user_id ON playlists(user_id);