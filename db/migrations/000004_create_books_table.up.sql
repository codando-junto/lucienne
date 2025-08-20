CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    author_id INTEGER NOT NULL REFERENCES authors (id)
);
