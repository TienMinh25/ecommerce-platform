drop index if exists idx_category_parent_id;
DROP TRIGGER IF EXISTS set_timestamp_categories ON categories;
DROP FUNCTION IF EXISTS update_modified_column;

drop table if exists categories;