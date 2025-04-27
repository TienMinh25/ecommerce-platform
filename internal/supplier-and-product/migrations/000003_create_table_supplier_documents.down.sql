alter table supplier_documents
drop constraint if exists verification_status_supplier_documents;

alter table supplier_documents
drop constraint if exists fk_supplier_id_supplier_documents;

drop trigger if exists set_timestamp_supplier_documents
on supplier_documents;

drop index if exists idx_supplier_id_supplier_documents;

drop table if exists supplier_documents;
