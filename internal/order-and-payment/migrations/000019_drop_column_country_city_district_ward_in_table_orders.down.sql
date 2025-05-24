alter table orders
add column country varchar(2000) not null,
add column city varchar(2000) not null,
add column district varchar(2000) not null,
add column ward varchar(2000);