alter table product_reviews
drop constraint if exists fk_product_id_product_reviews;

alter table product_reviews
drop constraint if exists valid_rating_value_product_reviews;

drop trigger if exists set_timestamp_product_reviews
on product_reviews;

drop trigger if exists trig_update_product_rating_on_insert
on product_reviews;
drop trigger if exists trig_update_product_rating_on_update
on product_reviews;
drop trigger if exists trig_update_product_rating_on_delete
on product_reviews;

drop function if exists update_product_rating_on_insert();
drop function if exists update_product_rating_on_update();
drop function if exists update_product_rating_on_delete();

drop index if exists idx_product_id_product_reviews;
drop index if exists idx_user_id_product_reviews;
drop index if exists idx_rating_product_reviews;