create table if not exists product_variants (
    id uuid primary key default gen_random_uuid(),
    product_id uuid not null,
    sku varchar(255) not null unique,
    variant_name varchar(2000) not null,
    price numeric(14, 2) not null,
    discount_price numeric(14, 2),
    inventory_quantity int not null default 0,
    low_stock_threshold int default 5,
    is_default boolean default false, -- use for display
    is_active boolean default true,
    shipping_class varchar(255) not null,
    image_url   varchar(3000) not null,
    alt_text    varchar(2500) not null,
    currency    varchar(20) not null default 'VND',
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

alter table product_variants
add constraint fk_product_id_product_variants
    foreign key (product_id) references products(id)
    on delete cascade;

alter table product_variants
add constraint check_price_non_negative
check (price >= 0);

alter table product_variants
add constraint check_discount_price_non_negative
check (discount_price is null or (discount_price >= 0 and discount_price <= price));

create index idx_product_id_product_variants
on product_variants (product_id);

create trigger set_timestamp_product_variants
    before update on product_variants
    for each row
    execute function update_modified_column();