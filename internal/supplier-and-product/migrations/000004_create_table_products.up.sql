create table if not exists products (
    id uuid primary key default gen_random_uuid(),
    supplier_id bigint not null,
    category_id bigint not null,
    name varchar(3000) not null,
    description text,
    image_url varchar(2000) not null,
    status varchar(20) not null,
    featured boolean default false,
    tax_class varchar(100) not null,
    sku_prefix varchar(100),
    average_rating decimal(3,2) default 0,
    total_reviews int default 0,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

alter table products
add constraint fk_supplier_id_products
    foreign key (supplier_id) references supplier_profiles (id)
        on delete cascade;

alter table products
add constraint fk_category_id_products
    foreign key (category_id) references categories (id)
    on delete set null;

create index if not exists idx_supplier_id_products
on products (supplier_id);

create index if not exists idx_category_id_products
on products (category_id);

create trigger set_timestamp_products
    before update on products
    for each row
    execute function update_modified_column();