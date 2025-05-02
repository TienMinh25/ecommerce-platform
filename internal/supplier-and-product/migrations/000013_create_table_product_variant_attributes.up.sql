create table if not exists product_variant_attributes (
    id uuid primary key default gen_random_uuid(),
    product_variant_id uuid not null,
    attribute_definition_id integer not null,
    attribute_option_id integer not null,
    text_value text,
    number_value decimal(12,2),
    boolean_value boolean,
    date_value date,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

-- foreign key
alter table product_variant_attributes
add constraint fk_product_variant_id_product_variant_attributes
    foreign key (product_variant_id) references
        product_variants (id) on delete cascade;

alter table product_variant_attributes
add constraint fk_attribute_definition_id_product_variant_attributes
    foreign key (attribute_definition_id) references
        attribute_definitions (id) on delete set null;

alter table product_variant_attributes
add constraint fk_attribute_option_id_product_variant_attributes
foreign key (attribute_option_id) references
    attribute_options(id) on delete set null;

-- index

create index if not exists idx_attribute_option_id_product_variant_attributes
on product_variant_attributes(attribute_option_id);

create index if not exists idx_product_variant_id_product_variant_attributes
on product_variant_attributes(product_variant_id);

create index if not exists idx_attribute_definition_id_product_variant_attributes
    on product_variant_attributes(attribute_definition_id);

-- bind trigger
create trigger set_timestamp_product_variant_attributes
    before update on product_variant_attributes
    for each row
    execute function update_modified_column();