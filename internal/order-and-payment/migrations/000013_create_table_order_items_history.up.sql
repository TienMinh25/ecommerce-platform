create table if not exists order_items_history (
    id uuid primary key default gen_random_uuid(),
    order_item_id uuid not null,
    status varchar(30) not null,
    notes text,
    created_at timestamptz default current_timestamp,
    created_by bigint not null
);

alter table order_items_history
add constraint fk_order_item_id_order_items_history
foreign key (order_item_id) references order_items(id) on delete cascade;

-- create index
create index idx_order_item_id_order_items_history
on order_items_history(order_item_id);

create index idx_created_by_order_items_history
on order_items_history(created_by);