DROP TRIGGER IF EXISTS set_timestamp_notifications ON notifications;
DROP FUNCTION IF EXISTS update_modified_column;

DROP INDEX IF EXISTS idx_noti_user_id;
DROP TABLE IF EXISTS notifications;