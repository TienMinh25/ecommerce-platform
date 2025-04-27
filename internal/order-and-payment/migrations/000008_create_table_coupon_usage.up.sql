create table if not exists coupon_usage (
    id uuid primary key default gen_random_uuid(),
    coupon_id uuid not null,
    user_id bigint not null,
    order_id uuid not null,
    discount_amount numeric(14,2) not null,
    used_at timestamptz
);

alter table coupon_usage
add constraint fk_coupon_id_coupon_usage
foreign key (coupon_id) references coupons(id) on delete no action;

alter table coupon_usage
add constraint fk_order_id_coupon_usage
foreign key (order_id) references orders(id) on delete no action;

alter table coupon_usage
add constraint check_discount_amount_coupon_usage
check ( discount_amount >= 0 );

alter table coupon_usage
add constraint unique_coupon_id_and_user_id_coupon_usage
unique (coupon_id, user_id);
