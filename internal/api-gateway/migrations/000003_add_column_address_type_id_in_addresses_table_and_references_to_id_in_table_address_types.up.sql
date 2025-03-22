ALTER TABLE addresses
ADD COLUMN address_type_id INTEGER;

ALTER TABLE addresses
ADD CONSTRAINT fk_addresses_address_types
FOREIGN KEY (address_type_id) REFERENCES address_types(id);
