alter table supplier_documents
add document_url varchar(2000);

alter table supplier_documents
drop column if exists documents jsonb;