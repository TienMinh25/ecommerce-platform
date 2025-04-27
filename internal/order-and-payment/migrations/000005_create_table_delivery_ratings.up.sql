create table if not exists delivery_ratings (
    id uuid primary key default gen_random_uuid(),
    delivery_person_id  bigint not null,
    order_id uuid not null,
    user_id bigint not null,
    rating smallint not null,
    comment text,
    created_at timestamptz
);

-- constraint
alter table delivery_ratings
add constraint fk_delivery_person_id_delivery_ratings
foreign key (delivery_person_id) references delivery_persons(id) on delete cascade;

alter table delivery_ratings
add constraint fk_order_id_delivery_ratings
foreign key (order_id) references orders(id) on delete set null;

alter table delivery_ratings
add constraint check_rating_delivery_ratings
check (rating >= 1 and rating <= 5);

-- create index
create index idx_delivery_person_id_delivery_ratings
on delivery_ratings(delivery_person_id);

create index idx_order_id_delivery_ratings
on delivery_ratings(order_id);

create index idx_user_id_delivery_ratings
on delivery_ratings(user_id);