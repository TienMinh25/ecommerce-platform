alter table user_payment_methods
drop constraint if exists fk_payment_method_id_user_payment_methods;

drop index if exists idx_user_id_user_payment_methods;
drop index if exists idx_payment_method_id_user_payment_methods;

drop trigger if exists set_timestamp_user_payment_methods
on user_payment_methods;

drop table if exists user_payment_methods;