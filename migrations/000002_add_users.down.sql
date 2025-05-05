-- Drop indexes
DROP INDEX IF EXISTS idx_mood_entries_user_id;
DROP INDEX IF EXISTS idx_playlists_user_id;

-- Remove user_id from mood_entries table
ALTER TABLE mood_entries DROP COLUMN IF EXISTS user_id;

-- Remove user_id from playlists table
ALTER TABLE playlists DROP COLUMN IF EXISTS user_id;

-- Drop users table
DROP TABLE IF EXISTS users;