CREATE TABLE company (
    id SERIAL PRIMARY KEY,
    owner_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(50),
    address TEXT,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_company_email ON company(email);

CREATE INDEX idx_company_name ON company(name);

CREATE TABLE car (
    id SERIAL PRIMARY KEY,
    company_id INT NOT NULL REFERENCES company(id) ON DELETE CASCADE,
    make VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INT NOT NULL,
    color VARCHAR(50),
    registration_no VARCHAR(50) UNIQUE NOT NULL,
    price_per_day DECIMAL(10, 2) NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX idx_car_company_id ON car(company_id);

CREATE INDEX idx_car_year_model ON car(year, model);

CREATE INDEX idx_car_make_model ON car(make, model);

CREATE UNIQUE INDEX idx_car_registration_no ON car(registration_no);

CREATE INDEX idx_car_price_per_day ON car(price_per_day);

