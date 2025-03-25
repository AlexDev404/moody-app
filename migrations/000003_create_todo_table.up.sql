CREATE TABLE IF NOT EXISTS todo (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    task text NOT NULL,
    completed boolean NOT NULL DEFAULT false
);
