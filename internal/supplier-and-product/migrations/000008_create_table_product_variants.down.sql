alter table product_variants
drop constraint if exists fk_product_id_product_variants;

alter table product_variants
drop constraint if exists check_price_non_negative;

alter table product_variants
drop constraint if exists check_discount_price_non_negative;

drop index if exists idx_product_id_product_variants;

drop trigger if exists set_timestamp_product_variants
on product_variants;

drop table if exists product_variants;