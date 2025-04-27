create table if not exists order_deliverers (
    id uuid primary key default gen_random_uuid(),
    order_item_id uuid not null unique,
    deliverer_id bigint not null,
    status varchar(50) not null,
    pickup_time timestamptz,
    delivery_time timestamptz,
    delivery_notes text,
    failure_reason text,
    proof_of_delivery varchar(2000),
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

CREATE TRIGGER set_timestamp_order_deliverers
    BEFORE UPDATE ON order_deliverers
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

-- foreign key
alter table order_deliverers
add constraint fk_deliverer_id_order_deliverers
foreign key (deliverer_id) references delivery_persons(id) on delete cascade;

alter table order_deliverers
add constraint check_status_order_deliverers
check (status in ('assigned', 'picked_up', 'in_transit', 'delivered', 'failed'));

-- create index
create index idx_deliverer_id_order_deliverers
on order_deliverers(deliverer_id);