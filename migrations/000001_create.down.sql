-- Drop indexes
DROP INDEX IF EXISTS idx_tracks_playlist_id;
DROP INDEX IF EXISTS idx_playlists_mood_id;
DROP INDEX IF EXISTS idx_mood_entries_created_at;

-- Drop tables
DROP TABLE IF EXISTS tracks;
DROP TABLE IF EXISTS playlists;
DROP TABLE IF EXISTS mood_entries;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";