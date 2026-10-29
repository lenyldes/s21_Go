CREATE TABLE IF NOT EXISTS address (
    id UUID PRIMARY KEY,
    country TEXT NOT NULL,
    city TEXT NOT NULL,
    street TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS client (
    id UUID PRIMARY KEY,
    client_name TEXT NOT NULL,
    client_surname TEXT NOT NULL,
    birthday DATE NOT NULL,
    gender TEXT NOT NULL,
    registration_date TIMESTAMPTZ NOT NULL,
    address_id UUID NOT NULL REFERENCES address(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS supplier (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    address_id UUID NOT NULL REFERENCES address(id) ON DELETE RESTRICT,
    phone_number TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY,
    image BYTEA NOT NULL
);

CREATE TABLE IF NOT EXISTS product (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    price NUMERIC NOT NULL CHECK (price >= 0),
    available_stock INTEGER NOT NULL CHECK (available_stock >= 0),
    last_update_date TIMESTAMPTZ NOT NULL,
    supplier_id UUID NOT NULL REFERENCES supplier(id) ON DELETE RESTRICT,
    image_id UUID REFERENCES images(id) ON DELETE SET NULL
);
