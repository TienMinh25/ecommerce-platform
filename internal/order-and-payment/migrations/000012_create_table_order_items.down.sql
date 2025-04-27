alter table order_deliverers
drop constraint if exists fk_order_item_id_order_deliverers;

alter table order_items
drop constraint if exists fk_order_id_order_items;

alter table order_items
drop constraint if exists check_status_order_items;

drop index if exists idx_order_id_order_items;

drop table if exists order_items;