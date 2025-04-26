alter table products_tags
drop constraint if exists fk_product_id_products_tags;

alter table products_tags
drop constraint if exists fk_tag_id_products_tags;

drop index if exists idx_product_id_products_tags;
drop index if exists idx_tag_id_products_tags;

drop table if exists products_tags;