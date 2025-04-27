create table if not exists supplier_documents (
    id uuid primary key default gen_random_uuid(),
    supplier_id bigint not null,
    document_url varchar(2000) not null,
    verification_status varchar(20) not null default 'pending',
    admin_note TEXT,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

alter table supplier_documents
add constraint verification_status_supplier_documents
    check (verification_status in ('pending', 'approved', 'rejected'));

alter table supplier_documents
add constraint fk_supplier_id_supplier_documents
    foreign key (supplier_id) references supplier_profiles (id)
        on delete cascade;

create index if not exists idx_supplier_id_supplier_documents
on supplier_documents (supplier_id);

create trigger set_timestamp_supplier_documents
    before update on supplier_documents
    for each row
    execute function update_modified_column();