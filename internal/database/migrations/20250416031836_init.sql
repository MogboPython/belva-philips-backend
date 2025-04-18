-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS public.users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    company_name TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    preferred_mode_of_communication TEXT NOT NULL,
    want_to_receive_text BOOLEAN DEFAULT FALSE,
    membership_status TEXT DEFAULT 'PAYG',
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public.orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    product_name TEXT NOT NULL,
    product_description TEXT,
    product_description_image TEXT,
    shoot_type TEXT NOT NULL,
    finish_type TEXT,
    quantity INTEGER NOT NULL,
    details JSONB,
    shots TEXT[],
    delivery_speed TEXT DEFAULT 'Standard',
    status TEXT DEFAULT 'QUOTE RECEIVED',
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES public.users (id) ON UPDATE NO ACTION ON DELETE NO ACTION
);


-- +goose Down
DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS orders;