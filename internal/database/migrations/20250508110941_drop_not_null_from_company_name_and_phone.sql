-- +goose Up
ALTER TABLE public.users 
ALTER COLUMN company_name DROP NOT NULL,
ALTER COLUMN phone_number DROP NOT NULL;

-- +goose Down
ALTER TABLE public.users 
ALTER COLUMN company_name SET NOT NULL,
ALTER COLUMN phone_number SET NOT NULL;
