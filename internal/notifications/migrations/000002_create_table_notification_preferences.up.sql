CREATE TABLE notification_preferences (
      user_id BIGINT PRIMARY KEY,
      email_preferences JSONB NOT NULL DEFAULT '{"order_status": true, "payment_status": true, "product_status": true, "promotion": true, "survey": true}',
      in_app_preferences JSONB NOT NULL DEFAULT '{"order_status": true, "payment_status": true, "product_status": true, "promotion": true, "survey": true}',
      created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER set_timestamp_notification_preferences
    BEFORE UPDATE ON notification_preferences
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();