-- drop constraint
alter table attribute_definitions
drop constraint if exists check_input_type_attribute_definitions;

-- drop trigger
drop trigger if exists set_timestamp_attribute_definitions
on attribute_definitions;

drop table if exists attribute_definitions;