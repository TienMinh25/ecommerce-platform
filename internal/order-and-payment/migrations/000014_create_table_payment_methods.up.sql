create table if not exists payment_methods (
    id serial primary key,
    name varchar(500) not null,
    code varchar(100) not null unique,
    is_active boolean default true,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

CREATE TRIGGER set_timestamp_payment_methods
    BEFORE UPDATE ON payment_methods
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();