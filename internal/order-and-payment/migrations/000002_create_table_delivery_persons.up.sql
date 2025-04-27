create table if not exists delivery_persons (
    id bigserial primary key,
    user_id bigint not null,
    id_card_number varchar(50) not null,
    vehicle_type varchar(100) not null,
    vehicle_license_plate varchar(50) not null,
    status varchar(20),
    average_rating decimal(3,2) not null default 0,
    total_rating int not null default 0,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

CREATE TRIGGER set_timestamp_delivery_persons
    BEFORE UPDATE ON delivery_persons
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

create index idx_user_id_delivery_persons
on delivery_persons(user_id);

alter table delivery_persons
add constraint unique_id_card_number_delivery_persons
unique (id_card_number);

alter table delivery_persons
add constraint unique_vehicle_license_plate_delivery_persons
unique (vehicle_license_plate);

alter table delivery_persons
add constraint check_status_delivery_persons
check (status in ('active', 'inactive', 'suspended'));