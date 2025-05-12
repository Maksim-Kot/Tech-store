CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE items (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL,
    price       NUMERIC(10, 2) NOT NULL,
    quantity    INTEGER NOT NULL,
    image_url   TEXT NOT NULL,
    attributes  JSONB NOT NULL,
    category_id BIGINT NOT NULL REFERENCES categories(id)
);