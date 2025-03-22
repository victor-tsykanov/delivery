-- +goose Up
-- +goose StatementBegin
create table orders
(
    id         uuid not null primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    location_x bigint,
    location_y bigint,
    status     varchar(100),
    courier_id uuid
);

create index idx_orders_deleted_at on orders (deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists orders;
-- +goose StatementEnd
