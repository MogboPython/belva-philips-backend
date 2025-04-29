-- +goose Up
ALTER TABLE public.orders 
DROP COLUMN product_description_image;

-- +goose Down
ALTER TABLE public.orders
ADD COLUMN product_description_image TEXT;
