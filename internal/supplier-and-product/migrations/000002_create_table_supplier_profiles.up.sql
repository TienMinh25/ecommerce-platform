create table if not exists supplier_profiles (
    id bigserial primary key,
    user_id bigint not null,
    company_name varchar(5000) not null,
    contact_phone varchar(20) not null,
    description text,
    logo_url varchar(2000) not null,
    business_address_id bigint not null,
    tax_id varchar(100) not null unique,
    status varchar(20) not null,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp,
    constraint status_supplier_profiles check(
        status in ('pending', 'active', 'suspended')
    )
);

create index if not exists idx_user_id_supplier_profiles
on supplier_profiles (user_id);

create index if not exists idx_company_name_supplier_profiles
on supplier_profiles (company_name);

create index if not exists idx_status_supplier_profiles
on supplier_profiles (status);

CREATE TRIGGER set_timestamp_supplier_profiles
    BEFORE UPDATE ON supplier_profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();