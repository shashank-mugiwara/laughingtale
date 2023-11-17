CREATE TABLE IF NOT EXISTS laughingtale.source_configs (
    identifier VARCHAR PRIMARY KEY,
    source_config JSONB DEFAULT '[]'::jsonb NOT NULL,
    created_at TIMESTAMP default NOW(),
    updated_at TIMESTAMP default NOW(),
    deleted_at TIMESTAMP default NOW()
);
