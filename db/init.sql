CREATE TABLE IF NOT EXISTS clients (
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name       TEXT NOT NULL,
    email      TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS documents (
    id                BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    client_id         BIGINT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    filename          TEXT NOT NULL,
    content_type      TEXT NOT NULL,
    size_bytes        BIGINT NOT NULL,
    storage_key       TEXT NOT NULL UNIQUE,
    status            TEXT NOT NULL DEFAULT 'pending'
                      CHECK (status IN ('pending', 'verified', 'rejected')),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- INFO: Since we don't have auth, the frontend uses this fixed 
-- user info to stand in for the logged-in user.
-- Do nothing on conflict because we need info to stay the same across table
-- creates, for takehome purposes only
INSERT INTO clients (name, email) VALUES
    ('John Smith, 'defaultuser@fake.com')
ON CONFLICT (email) DO NOTHING;
