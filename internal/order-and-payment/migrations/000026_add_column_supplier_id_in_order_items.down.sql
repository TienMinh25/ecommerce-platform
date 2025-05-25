drop index if exists idx_order_items_supplier_id on order_items;

alter table order_items
drop column supplier_id;