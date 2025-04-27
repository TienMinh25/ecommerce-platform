create table if not exists inventory_transactions (
    id uuid primary key default gen_random_uuid(),
    product_variant_id uuid not null,
    quantity_change int not null,
    previous_quantity int not null,
    new_quantity int not null,
    transaction_type varchar(200) not null,
    performed_by bigint not null,
    created_at timestamptz default current_timestamp
);

alter table inventory_transactions
add constraint fk_product_variant_id_inventory_transactions
    foreign key (product_variant_id) references
    product_variants (id) on delete cascade;

alter table inventory_transactions
add constraint check_transaction_type_inventory_transactions
check (transaction_type in ('purchase', 'sale', 'return', 'adjustment', 'inventory_count'));

create index idx_product_variant_id_inventory_transactions
on inventory_transactions(product_variant_id);

create index idx_transaction_type_inventory_transactions
on inventory_transactions(transaction_type);

create index idx_performed_by_inventory_transactions
on inventory_transactions(performed_by);

create index idx_created_at_inventory_transactions
on inventory_transactions(created_at);