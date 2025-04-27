create table if not exists products_tags (
    product_id uuid not null,
    tag_id uuid not null
);

create index if not exists idx_product_id_products_tags
on products_tags (product_id);

create index if not exists idx_tag_id_products_tags
on products_tags (tag_id);

alter table products_tags
add primary key (product_id, tag_id);

alter table products_tags
add constraint fk_product_id_products_tags foreign key (product_id)
    references products(id) on delete cascade;

alter table products_tags
add constraint fk_tag_id_products_tags foreign key (tag_id)
    references tags(id) on delete cascade;