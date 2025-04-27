create table if not exists attribute_options (
    id serial primary key,
    attribute_definition_id integer not null,
    option_value varchar(500) not null,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

alter table attribute_options
add constraint fk_attribute_definition_id_attribute_options
    foreign key (attribute_definition_id) references
    attribute_definitions(id) on delete cascade;

create index idx_attribute_definition_id_attribute_options
on attribute_options(attribute_definition_id);

-- add constraint unique
alter table attribute_options
add constraint unique_attribute_definition_id_and_option_value_attribute_options
    unique (attribute_definition_id, option_value);

-- bind trigger
create trigger set_timestamp_attribute_options
    before update on attribute_options
    for each row
    execute function update_modified_column();