-- Sessions for folio's SessionStore: the browser cookie holds only the id.
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    data JSONB NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
