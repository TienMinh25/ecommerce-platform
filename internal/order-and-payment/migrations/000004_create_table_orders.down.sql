drop trigger if exists set_timestamp_orders
on orders;

drop index if exists idx_user_id_orders;

drop table if exists orders;