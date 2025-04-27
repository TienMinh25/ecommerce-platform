create table if not exists carts (
    id bigserial primary key,
    user_id bigint not null
);

create index idx_user_id_carts
on carts (user_id);