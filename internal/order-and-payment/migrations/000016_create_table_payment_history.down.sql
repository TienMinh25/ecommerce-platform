alter table payment_history
drop constraint if exists fk_order_item_id_payment_history;

alter table payment_history
drop constraint if exists fk_user_payment_method_id_payment_history;

alter table payment_history
drop constraint if exists check_status_payment_history;

drop index if exists idx_order_item_id_payment_history;
drop index if exists idx_user_payment_method_id_payment_history;

drop table if exists payment_history;