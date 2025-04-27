create table if not exists order_items (
    id uuid primary key default gen_random_uuid(),
    order_id uuid not null,
    product_name varchar(2000) not null,
    product_sku  varchar(255) not null,
    product_variant_image_url varchar(2000) not null,
    product_variant_name varchar(2000) not null,
    quantity int not null,
    unit_price numeric(14, 2) not null,
    discount_price numeric(14, 2),
    total_price numeric(14, 2),
    attributes jsonb,
    estimated_delivery_date	date,
    actual_delivery_date date,
    cancelled_reason text,
    notes text,
    status varchar(50) not null,
    shipping_fee numeric(14, 2) not null,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

-- add constraint because my error thinking
-- just changed a little bit
alter table order_deliverers
add constraint fk_order_item_id_order_deliverers
foreign key (order_item_id) references order_items(id);

-- main part for order_items table
alter table order_items
add constraint fk_order_id_order_items
foreign key (order_id) references orders(id) on delete cascade;

alter table order_items
add constraint check_status_order_items
check ( status in ('pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled', 'refunded') );

-- create index
create index idx_order_id_order_items
on order_items(order_id);

