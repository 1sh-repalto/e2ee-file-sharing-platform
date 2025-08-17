CREATE TABLE shares (
    id UUID PRIMARY KEY,
    file_id UUID REFERENCES files(id) ON DELETE CASCADE,
    recipient_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    wrapped_key BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);