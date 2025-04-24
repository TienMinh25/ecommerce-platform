create table user_devices (
    id uuid primary key default gen_random_uuid(),
    user_id bigint not null,
    device_token text not null,
    platform varchar(50) not null,
    last_active_at timestamptz
);

create index if not exists idx_user_id_user_devices
on user_devices(user_id);