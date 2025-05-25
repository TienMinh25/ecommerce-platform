alter table order_items
drop constraint if exists check_status_order_items;

alter table order_items
    add constraint check_status_order_items
        check (status in (
                          'pending',
                          'confirmed',
                          'processing',
                          'shipped',
                          'delivered',
                          'cancelled',
                          'refunded'
            ));