alter table supplier_documents
drop column if exists document_url;

alter table supplier_documents
add documents jsonb;