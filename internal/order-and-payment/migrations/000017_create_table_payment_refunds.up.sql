create table if not exists payment_refunds (
    id uuid primary key default gen_random_uuid(),
    payment_history_id uuid not null,
    amount numeric(14, 2) not null,
    status varchar(50) not null,
    transaction_id varchar(255),
    payment_gateway_response JSONB,
    processed_by bigint not null,
    refunded_at	timestamptz
);

alter table payment_refunds
add constraint fk_payment_history_id_payment_refunds
foreign key (payment_history_id) references
payment_history(id);

alter table payment_refunds
add constraint check_status_payment_refunds
    CHECK (status IN ('pending', 'processing', 'completed', 'failed'));

create index idx_payment_history_id_payment_refunds
on payment_refunds(payment_history_id);

create index idx_processed_by_payment_refunds
on payment_refunds(processed_by);