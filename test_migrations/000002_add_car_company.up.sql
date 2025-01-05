-- Create the 'company' table
CREATE TABLE company (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    phone TEXT,
    address TEXT,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_company_email ON company(email);
CREATE INDEX idx_company_name ON company(name);

CREATE TABLE car (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    company_id INTEGER NOT NULL,
    make TEXT NOT NULL,
    model TEXT NOT NULL,
    year INTEGER NOT NULL,
    color TEXT,
    registration_no TEXT UNIQUE NOT NULL,
    price_per_day REAL NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES company(id) ON DELETE CASCADE
);

CREATE INDEX idx_car_company_id ON car(company_id);
CREATE INDEX idx_car_year_model ON car(year, model);
CREATE INDEX idx_car_make_model ON car(make, model);
CREATE INDEX idx_car_registration_no ON car(registration_no);
CREATE INDEX idx_car_price_per_day ON car(price_per_day);
