alter table coupon_usage
drop constraint if exists fk_coupon_id_coupon_usage;

alter table coupon_usage
drop constraint if exists fk_order_id_coupon_usage;

alter table coupon_usage
drop constraint if exists check_discount_amount_coupon_usage;

alter table coupon_usage
drop constraint if exists unique_coupon_id_and_user_id_coupon_usage;

drop table if exists coupon_usage;