create table if not exists product_reviews (
    id uuid primary key default gen_random_uuid(),
    product_id uuid not null,
    user_id bigint not null,
    rating smallint not null,
    comment TEXT,
    is_verified_purchase boolean default false,
    helpful_votes int default 0,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

alter table product_reviews
add constraint fk_product_id_product_reviews
    foreign key (product_id) references products (id)
    on delete cascade;

create index if not exists idx_product_id_product_reviews
on product_reviews (product_id);

create index if not exists idx_user_id_product_reviews
on product_reviews (user_id);

create index if not exists idx_rating_product_reviews
on product_reviews (rating);

alter table product_reviews
add constraint valid_rating_value_product_reviews
    check ( rating >= 1 and rating <= 5 );

create trigger set_timestamp_product_reviews
    before update on product_reviews
    for each row
    execute function update_modified_column();

-- Trigger cho trường hợp INSERT vào product_reviews
CREATE OR REPLACE FUNCTION update_product_rating_on_insert()
RETURNS TRIGGER AS $$
BEGIN
    -- Cập nhật average_rating và total_reviews trong products
UPDATE products
SET
    total_reviews = (SELECT COUNT(*) FROM product_reviews WHERE product_id = NEW.product_id),
    average_rating = (SELECT COALESCE(AVG(rating), 0) FROM product_reviews WHERE product_id = NEW.product_id)
WHERE id = NEW.product_id;

RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trig_update_product_rating_on_insert
    AFTER INSERT ON product_reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_product_rating_on_insert();

-- Trigger cho trường hợp UPDATE trong product_reviews
CREATE OR REPLACE FUNCTION update_product_rating_on_update()
RETURNS TRIGGER AS $$
BEGIN
    -- Cập nhật average_rating trong products
UPDATE products
SET
    average_rating = (SELECT COALESCE(AVG(rating), 0) FROM product_reviews WHERE product_id = NEW.product_id)
WHERE id = NEW.product_id;

RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trig_update_product_rating_on_update
    AFTER UPDATE ON product_reviews
    FOR EACH ROW
    WHEN (OLD.rating IS DISTINCT FROM NEW.rating)
EXECUTE FUNCTION update_product_rating_on_update();

-- Trigger cho trường hợp DELETE từ product_reviews
CREATE OR REPLACE FUNCTION update_product_rating_on_delete()
RETURNS TRIGGER AS $$
BEGIN
    -- Cập nhật average_rating và total_reviews trong products
UPDATE products
SET
    total_reviews = (SELECT COUNT(*) FROM product_reviews WHERE product_id = OLD.product_id),
    average_rating = (SELECT COALESCE(AVG(rating), 0) FROM product_reviews WHERE product_id = OLD.product_id)
WHERE id = OLD.product_id;

RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trig_update_product_rating_on_delete
    AFTER DELETE ON product_reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_product_rating_on_delete();