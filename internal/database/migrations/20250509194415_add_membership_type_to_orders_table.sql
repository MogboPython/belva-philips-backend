-- +goose Up
ALTER TABLE public.orders 
ADD COLUMN membership_type TEXT DEFAULT 'PAY AS YOU GO';

-- +goose Down
ALTER TABLE public.orders
DROP COLUMN membership_type;
