drop trigger if exists set_timestamp_order_deliverers
on order_deliverers;

alter table order_deliverers
drop constraint if exists fk_deliverer_id_order_deliverers;

alter table order_deliverers
drop constraint if exists check_status_order_deliverers;

drop index if exists idx_order_id_order_deliverers;

drop index if exists idx_deliverer_id_order_deliverers;

drop table if exists order_deliverers;