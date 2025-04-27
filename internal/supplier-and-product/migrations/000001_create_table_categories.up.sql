create table categories (
    id bigserial primary key,
    name varchar(2000) not null,
    description text,
    is_active boolean default false,
    deleted_at timestamptz,
    parent_id bigint,
    image_url text not null,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp,
    constraint fk_parent_category foreign key (parent_id)
        references categories(id) on delete set null
);

-- function trigger update timestamps when update
CREATE OR REPlACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp_categories
    BEFORE UPDATE ON categories
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

create index idx_category_parent_id
on categories(parent_id);
