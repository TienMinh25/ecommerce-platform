ALTER TABLE addresses 
ADD COLUMN address_type VARCHAR(50) NOT NULL;

DROP TRIGGER IF EXISTS set_timestamp_address_types
ON address_types;

DROP TABLE IF EXISTS address_types;