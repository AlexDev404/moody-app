CREATE TABLE IF NOT EXISTS journal (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    content text NOT NULL
);
