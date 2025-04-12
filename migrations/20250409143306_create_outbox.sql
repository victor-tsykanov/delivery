-- +goose Up
-- +goose StatementBegin
create table outbox
(
    id           text not null primary key,
    created_at   timestamp with time zone,
    updated_at   timestamp with time zone,
    deleted_at   timestamp with time zone,
    type         varchar(500),
    payload      bytea,
    processed_at timestamp with time zone
);
create index idx_outbox_processed_at on public.outbox (processed_at);
create index idx_outbox_deleted_at on public.outbox (deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table outbox;
-- +goose StatementEnd
