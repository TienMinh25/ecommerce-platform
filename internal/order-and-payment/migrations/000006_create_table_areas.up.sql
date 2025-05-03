create table if not exists areas (
    id bigserial primary key,
    city varchar(2000) not null,
    country varchar(2000) not null,
    district varchar(2000) not null,
    ward varchar(2000),
    area_code varchar(200) not null unique,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

CREATE TRIGGER set_timestamp_areas
    BEFORE UPDATE ON areas
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();
