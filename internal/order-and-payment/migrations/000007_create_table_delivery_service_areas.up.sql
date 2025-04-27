create table if not exists delivery_service_areas (
    id bigserial primary key,
    delivery_person_id bigint not null,
    area_id bigint not null,
    is_active boolean not null default true,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

alter table delivery_service_areas
add constraint fk_delivery_person_id_delivery_service_areas
    foreign key (delivery_person_id) references delivery_persons(id)
        on delete cascade;

alter table delivery_service_areas
add constraint fk_area_id_delivery_service_areas
    foreign key (area_id) references areas(id)
    on delete cascade;

CREATE TRIGGER set_timestamp_delivery_service_areas
    BEFORE UPDATE ON delivery_service_areas
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

create index idx_delivery_person_id_delivery_service_areas
on delivery_service_areas(delivery_person_id);

create index idx_area_id_delivery_service_areas
on delivery_service_areas(area_id);

