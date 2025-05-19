-- +goose Up
ALTER TABLE public.users 
DROP COLUMN preferred_mode_of_communication, 
DROP COLUMN want_to_receive_text;


-- +goose Down
ALTER TABLE public.users
ADD COLUMN preferred_mode_of_communication TEXT,
ADD COLUMN want_to_receive_text BOOLEAN DEFAULT FALSE;

-- Update existing rows to set preferred_mode_of_communication to a default value
-- since it was originally NOT NULL
UPDATE public.users SET preferred_mode_of_communication = 'email' WHERE preferred_mode_of_communication IS NULL;

-- Now add the NOT NULL constraint
ALTER TABLE public.users
ALTER COLUMN preferred_mode_of_communication SET NOT NULL;