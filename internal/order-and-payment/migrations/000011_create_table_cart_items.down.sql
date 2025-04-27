drop trigger if exists set_timestamp_cart_items
on cart_items;

alter table cart_items
drop constraint if exists fk_cart_id_cart_items;

drop index if exists idx_cart_id_cart_items;
drop index if exists idx_product_id_cart_items;
drop index if exists idx_product_variant_id_cart_items;

drop table if exists cart_items;