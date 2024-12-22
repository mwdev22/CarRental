CREATE TABLE users (
    id SERIAL PRIMARY KEY,                    
    username VARCHAR(20) NOT NULL,             
    email VARCHAR(40) NOT NULL UNIQUE,  
    role INT DEFAULT 1,       
    password BYTEA NOT NULL,  
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);

CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_users_email ON users (email);