CREATE TABLE IF NOT EXISTS role_user_permissions (
    role_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    permission_detail JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, user_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TRIGGER set_timestamp_role_user_permissions
BEFORE UPDATE ON role_user_permissions
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();