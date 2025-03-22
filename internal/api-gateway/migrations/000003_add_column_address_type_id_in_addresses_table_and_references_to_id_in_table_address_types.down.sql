ALTER TABLE addresses
DROP CONSTRAINT fk_addresses_address_types;

ALTER TABLE addresses
DROP COLUMN address_type_id;