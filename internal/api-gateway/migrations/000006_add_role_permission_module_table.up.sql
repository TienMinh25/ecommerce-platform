CREATE TABLE IF NOT EXISTS role_permission_module (
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    permission_detail TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id)
);

CREATE TRIGGER set_timestamp_role_permission_module
BEFORE UPDATE ON role_permission_module
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();
