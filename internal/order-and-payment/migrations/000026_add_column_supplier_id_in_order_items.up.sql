alter table order_items
add column supplier_id bigint not null;

create index idx_order_items_supplier_id
on order_items(supplier_id);