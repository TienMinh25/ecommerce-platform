create table if not exists orders (
    id uuid primary key default gen_random_uuid(),
    user_id bigint not null,
    tracking_number varchar(100) unique,
    shipping_address text not null,
    country varchar(1000) not null,
    city varchar(1000) not null,
    district varchar(1000) not null,
    ward varchar(1000),
    shipping_method  varchar(100) not null,
    sub_total numeric(14,2) not null,
    discount_amount numeric(14, 2),
    tax_amount  numeric(14, 2),
    total_amount numeric(14, 2) not null,
    recipient_name varchar(2000) not null,
    recipient_phone varchar(20) not null,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

create index idx_user_id_orders
on orders(user_id);

CREATE TRIGGER set_timestamp_orders
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

create index idx_recipient_name_orders
on orders(recipient_name);

create index idx_recipient_phone_orders
on orders(recipient_phone);