-- +goose Up
ALTER TABLE public.orders 
ADD COLUMN order_name TEXT NOT NULL DEFAULT ''::text;

-- +goose Down
ALTER TABLE public.orders
DROP COLUMN order_name;
