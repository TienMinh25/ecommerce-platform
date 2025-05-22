alter table coupons
drop constraint if exists check_applies_to_coupons;

alter table coupons
drop column applies_to;