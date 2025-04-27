drop trigger if exists set_timestamp_delivery_persons
on delivery_persons;

drop index if exists idx_user_id_delivery_persons;

alter table delivery_persons
drop constraint if exists unique_id_card_number_delivery_persons;

alter table delivery_persons
drop constraint if exists unique_vehicle_license_plate_delivery_persons;

alter table delivery_persons
drop constraint if exists check_status_delivery_persons;

drop table if exists delivery_persons;