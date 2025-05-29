alter table supplier_documents
    add column document_type varchar(100) not null default 'register';