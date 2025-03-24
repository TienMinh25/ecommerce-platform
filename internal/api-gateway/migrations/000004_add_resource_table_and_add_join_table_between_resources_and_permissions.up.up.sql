CREATE TABLE IF NOT EXISTS resources(
                                        id SERIAL PRIMARY KEY,
                                        name VARCHAR(100) NOT NULL,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
CREATE TRIGGER set_timestamp_resources
    BEFORE UPDATE ON resources
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();