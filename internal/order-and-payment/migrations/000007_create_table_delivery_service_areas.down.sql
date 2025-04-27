alter table delivery_service_areas
drop constraint if exists fk_delivery_person_id_delivery_service_areas;

alter table delivery_service_areas
drop constraint if exists fk_area_id_delivery_service_areas;

drop trigger if exists set_timestamp_delivery_service_areas
on delivery_service_areas;

drop index if exists idx_delivery_person_id_delivery_service_areas;

drop index if exists idx_area_id_delivery_service_areas;

drop table if exists delivery_service_areas;