create table if not exists delivery_person_applications (
    id bigserial primary key,
    user_id bigint not null,
    id_card_number varchar(50) not null,
    id_card_front_image varchar(2000) not null,
    id_card_back_iamge varchar(2000) not null,
    vehicle_type varchar(100) not null,
    vehicle_license_plate varchar(50) not null,
    service_area    jsonb not null,
    application_status varchar(20) not null,
    rejection_reason text,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

CREATE TRIGGER set_timestamp_delivery_person_applications
    BEFORE UPDATE ON delivery_person_applications
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

-- create index
create index idx_user_id_delivery_person_applications
on delivery_person_applications(user_id);

create index idx_id_card_number_delivery_person_applications
on delivery_person_applications(id_card_number);

create index idx_vehicle_license_plate_delivery_person_applications
on delivery_person_applications(vehicle_license_plate);

alter table delivery_person_applications
add constraint check_application_status_delivery_person_applications
check ( application_status in ('pending', 'approved', 'rejected') );