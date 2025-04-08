CREATE TABLE IF NOT EXISTS role_permissions (
    role_id BIGINT NOT NULL PRIMARY KEY,
    permission_detail JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

CREATE TRIGGER set_timestamp_role_permissions
BEFORE UPDATE ON role_permissions
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();