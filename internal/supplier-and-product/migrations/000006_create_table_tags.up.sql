create table if not exists tags (
    id uuid primary key default gen_random_uuid(),
    name varchar(400) not null unique,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

create trigger set_timestamp_tags
    before update on tags
    for each row
    execute function update_modified_column();