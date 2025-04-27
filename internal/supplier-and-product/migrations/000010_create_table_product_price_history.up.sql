create table if not exists product_price_history (
    id uuid primary key default gen_random_uuid(),
    product_variant_id uuid not null,
    old_price numeric(14, 2) not null,
    new_price numeric(14, 2) not null,
    old_discount_price numeric(14, 2),
    new_discount_price numeric(14, 2),
    changed_by bigint not null,
    reason text,
    created_at timestamptz default current_timestamp
);

alter table product_price_history
add constraint fk_product_variant_id_product_price_history
foreign key (product_variant_id) references
    product_variants(id) on delete cascade;

alter table product_price_history
add constraint check_old_price_non_negative_product_price_history
check (old_price >= 0);

alter table product_price_history
add constraint check_new_price_non_negative_product_price_history
check (new_price >= 0);

alter table product_price_history
add constraint check_old_discount_price_non_negative_product_price_history
check (old_discount_price is null or (old_discount_price >= 0 and old_discount_price <= old_price));

alter table product_price_history
add constraint check_new_discount_price_non_negative_product_price_history
check (new_discount_price is null or (new_discount_price >= 0 and new_discount_price <= new_price));

create index idx_product_variant_id_product_price_history
on product_price_history(product_variant_id);

create index idx_changed_by_product_price_history
on product_price_history(changed_by);

create index idx_created_at_product_price_history
on product_price_history(created_at);