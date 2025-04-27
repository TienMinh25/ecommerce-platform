create table if not exists payment_history (
    id uuid primary key default gen_random_uuid(),
    order_item_id uuid not null unique,
    user_payment_method_id bigint not null,
    amount numeric(14, 2) not null,
    currency varchar(20) not null default 'VND',
    status varchar(50) not null,
    transaction_id varchar(255),
    payment_gateway varchar(100),
    payment_gateway_response jsonb,
    error_message text,
    paid_at timestamptz,
    created_at timestamptz default current_timestamp
);

alter table payment_history
add constraint fk_order_item_id_payment_history
foreign key (order_item_id) references order_items(id)
on delete no action;

alter table payment_history
add constraint fk_user_payment_method_id_payment_history
foreign key (user_payment_method_id) references user_payment_methods(id)
on delete no action;

alter table payment_history
add constraint check_status_payment_history
check ( status in ('pending', 'processing', 'completed', 'failed', 'refunded', 'partially_refunded') );

create index idx_order_item_id_payment_history
on payment_history(order_item_id);

create index idx_user_payment_method_id_payment_history
on payment_history(user_payment_method_id);