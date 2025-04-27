alter table payment_refunds
drop constraint if exists fk_payment_history_id_payment_refunds;

alter table payment_refunds
drop constraint if exists check_status_payment_refunds;

drop index if exists idx_payment_history_id_payment_refunds;
drop index if exists idx_processed_by_payment_refunds;

drop table if exists payment_refunds;