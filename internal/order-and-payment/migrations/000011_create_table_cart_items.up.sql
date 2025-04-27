create table if not exists cart_items (
    id uuid primary key default gen_random_uuid(),
    cart_id bigint not null,
    product_id uuid not null,
    quantity int not null,
    product_variant_id uuid not null,
    added_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

CREATE TRIGGER set_timestamp_cart_items
    BEFORE UPDATE ON cart_items
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

alter table cart_items
add constraint fk_cart_id_cart_items
foreign key (cart_id) references carts(id)
on delete cascade;

create index idx_cart_id_cart_items
on cart_items(cart_id);

create index idx_product_id_cart_items
on cart_items(product_id);

create index idx_product_variant_id_cart_items
on cart_items(product_variant_id);