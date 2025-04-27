CREATE INDEX IF NOT EXISTS idx_users_roles_user_id
ON users_roles (user_id);

CREATE INDEX IF NOT EXISTS idx_users_roles_role_id
ON users_roles (role_id);