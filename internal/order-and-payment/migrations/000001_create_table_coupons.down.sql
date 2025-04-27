alter table coupons
drop constraint check_discount_type_coupons;

alter table coupons
drop constraint check_applies_to_coupons;

drop index if exists idx_start_date_coupons;
drop index if exists idx_end_date_coupons;

DROP TRIGGER IF EXISTS set_timestamp_coupons ON coupons;
DROP FUNCTION IF EXISTS update_modified_column();

drop table if exists coupons;