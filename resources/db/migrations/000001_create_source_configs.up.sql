CREATE TABLE IF NOT EXISTS laughingtale.source_configs (
    identifier VARCHAR PRIMARY KEY,
    source_config JSONB DEFAULT '[]'::jsonb NOT NULL,
    poller_config JSONB DEFAULT '{}'::jsonb NOT NULL,
    type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP default NOW(),
    updated_at TIMESTAMP default NOW()
);
