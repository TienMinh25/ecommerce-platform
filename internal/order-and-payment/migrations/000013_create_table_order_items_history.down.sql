alter table order_items_history
drop constraint if exists fk_order_item_id_order_items_history;

drop index if exists idx_order_item_id_order_items_history;
drop index if exists idx_created_by_order_items_history;

drop table if exists order_items_history;