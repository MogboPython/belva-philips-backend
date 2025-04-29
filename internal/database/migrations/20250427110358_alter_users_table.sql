-- +goose Up
ALTER TABLE public.users 
DROP COLUMN preferred_mode, 
DROP COLUMN want_to_receive_text;


-- +goose Down
ALTER TABLE public.users
ADD COLUMN preferred_mode TEXT,
ADD COLUMN want_to_receive_text BOOLEAN DEFAULT FALSE;

-- Update existing rows to set preferred_mode to a default value
-- since it was originally NOT NULL
UPDATE public.users SET preferred_mode = 'email' WHERE preferred_mode IS NULL;

-- Now add the NOT NULL constraint
ALTER TABLE public.users
ALTER COLUMN preferred_mode SET NOT NULL;