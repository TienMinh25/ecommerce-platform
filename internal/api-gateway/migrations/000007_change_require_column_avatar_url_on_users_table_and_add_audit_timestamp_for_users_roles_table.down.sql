ALTER TABLE users_roles
DROP COLUMN created_at,
DROP COLUMN updated_at;

DROP TRIGGER IF EXISTS users_roles_updated_at_trigger ON users_roles;