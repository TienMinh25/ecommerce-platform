alter table supplier_profiles
drop constraint status_supplier_profiles;

drop trigger if exists set_timestamp_supplier_profiles
on supplier_profiles;

drop index if exists idx_user_id_supplier_profiles;
drop index if exists idx_company_name_supplier_profiles;
drop index if exists idx_status_supplier_profiles;

drop table if exists supplier_profiles;