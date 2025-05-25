alter table order_items
drop constraint if exists check_status_order_items;

alter table order_items
add constraint check_status_order_items
check (status in (
    'pending_payment',
    'pending',
    'confirmed',
    'processing',
    'ready_to_ship',
    'in_transit',
    'out_for_delivery',
    'delivered',
    'cancelled',
    'payment_failed',
    'refunded'
    ));