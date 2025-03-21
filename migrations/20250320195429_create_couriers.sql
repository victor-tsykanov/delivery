-- +goose Up
-- +goose StatementBegin
create table if not exists couriers
(
    id         uuid not null primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name       varchar(500),
    location_x bigint,
    location_y bigint,
    status     varchar(100)
);

create index if not exists idx_couriers_deleted_at on couriers (deleted_at);

create table if not exists transports
(
    id         uuid not null primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    courier_id uuid constraint fk_couriers_transport references couriers,
    name       varchar(500),
    speed      bigint
);

create index if not exists idx_transports_deleted_at on transports (deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists transports;
drop table if exists couriers;
-- +goose StatementEnd
