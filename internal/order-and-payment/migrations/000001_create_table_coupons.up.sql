create table if not exists coupons (
    id uuid primary key default gen_random_uuid(),
    code varchar(100) not null unique,
    name varchar(5000) not null,
    description TEXT,
    discount_type varchar(20) not null,
    discount_value  numeric(14,2) not null default 0,
    maximum_discount_amount numeric(14, 2) not null default 0,
    minimum_order_amount numeric(14, 2) not null default 0,
    currency VARCHAR(20) NOT NULL DEFAULT 'VND',
    start_date  timestamptz not null,
    end_date    timestamptz not null,
    usage_limit int,
    usage_count int default 0,
    is_active   boolean default true,
    applies_to  varchar(20) not null,
    created_at  timestamptz default current_timestamp,
    updated_at  timestamptz default current_timestamp
);

-- function trigger update timestamps when update
CREATE OR REPlACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp_coupons
    BEFORE UPDATE ON coupons
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

alter table coupons
add constraint check_discount_type_coupons
    check ( discount_type in ('percentage', 'fixed_amount') );

alter table coupons
add constraint check_applies_to_coupons
check (applies_to in ('order', 'product', 'category'));

create index idx_start_date_coupons
on coupons(start_date);

create index idx_end_date_coupons
on coupons(start_date);