alter table products
drop constraint if exists fk_supplier_id_products;

alter table products
drop constraint if exists fk_category_id_products;

drop trigger if exists set_timestamp_products
on products;

drop index if exists idx_supplier_id_products;
drop index if exists idx_category_id_products;

drop table if exists products;