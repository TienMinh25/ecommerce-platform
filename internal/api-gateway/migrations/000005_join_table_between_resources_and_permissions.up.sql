CREATE TABLE IF NOT EXISTS permissions_resources(
    resource_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INTEGER REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (resource_id, permission_id)
);