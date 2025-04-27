create table if not exists user_payment_methods (
    id bigserial primary key,
    user_id bigint not null,
    payment_method_id integer not null,
    is_default boolean default false,
    card_holder_name varchar(255) not null,
    card_number varchar(255) not null,
    card_expiry_month smallint,
    card_expiry_year smallint,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

alter table user_payment_methods
add constraint fk_payment_method_id_user_payment_methods
foreign key (payment_method_id) references payment_methods(id)
on delete cascade;

create index idx_user_id_user_payment_methods
on user_payment_methods(user_id);

create index idx_payment_method_id_user_payment_methods
    on user_payment_methods(payment_method_id);

CREATE TRIGGER set_timestamp_user_payment_methods
    BEFORE UPDATE ON user_payment_methods
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();