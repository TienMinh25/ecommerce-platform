alter table delivery_ratings
drop constraint if exists fk_delivery_person_id_delivery_ratings;

alter table delivery_ratings
drop constraint if exists fk_order_id_delivery_ratings;

alter table delivery_ratings
drop constraint if exists check_rating_delivery_ratings;

drop index if exists idx_delivery_person_id_delivery_ratings;
drop index if exists idx_order_id_delivery_ratings;
drop index if exists idx_user_id_delivery_ratings;

drop table if exists delivery_ratings;