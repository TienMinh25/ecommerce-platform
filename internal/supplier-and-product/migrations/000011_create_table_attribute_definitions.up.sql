create table if not exists attribute_definitions (
    id serial primary key,
    name varchar(500) not null unique,
    description text,
    is_filterable boolean default true,
    is_required boolean default false,
    input_type varchar(100) not null,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

-- add constraint check for input_type
alter table attribute_definitions
add constraint check_input_type_attribute_definitions
    check (input_type in ('text', 'number', 'select', 'multiselect', 'boolean', 'date'));

-- bind trigger to updated_at
create trigger set_timestamp_attribute_definitions
    before update on attribute_definitions
    for each row
    execute function update_modified_column();