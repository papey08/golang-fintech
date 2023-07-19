CREATE TABLE ads (
    id SERIAL PRIMARY KEY NOT NULL,
    title VARCHAR(100),
    text VARCHAR(500),
    author_id INTEGER,
    published BOOLEAN,
    creation_date DATE,
    update_date DATE
);