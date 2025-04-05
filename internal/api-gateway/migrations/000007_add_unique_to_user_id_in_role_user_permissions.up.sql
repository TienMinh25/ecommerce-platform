ALTER TABLE role_user_permissions
    ADD CONSTRAINT role_user_permissions_user_id_unique UNIQUE (user_id);
