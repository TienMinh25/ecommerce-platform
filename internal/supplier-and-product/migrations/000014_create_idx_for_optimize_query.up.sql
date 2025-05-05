CREATE INDEX idx_products_name_gin ON products USING gin(to_tsvector('simple', name));

-- Index cho sắp xếp
CREATE INDEX idx_products_rating_reviews ON products(average_rating DESC, total_reviews DESC);