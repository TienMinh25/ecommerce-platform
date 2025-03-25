-- Drop triggers
DROP TRIGGER IF EXISTS set_timestamp_users ON users;
DROP TRIGGER IF EXISTS set_timestamp_roles ON roles;
DROP TRIGGER IF EXISTS set_timestamp_user_password ON user_password;
DROP TRIGGER IF EXISTS set_timestamp_addresses ON addresses;

-- Drop functions
DROP FUNCTION IF EXISTS update_permission_name;
DROP FUNCTION IF EXISTS update_modified_column;

-- Drop indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_refresh_token_token;
DROP INDEX IF EXISTS idx_refresh_token_user_id;
DROP INDEX IF EXISTS idx_addresses_user_id;

-- Drop tables in reverse order to avoid foreign key conflicts
DROP TABLE IF EXISTS refresh_token;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS user_password;
DROP TABLE IF EXISTS users_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;
