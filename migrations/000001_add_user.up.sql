CREATE TABLE users (
    id SERIAL PRIMARY KEY,                    
    username VARCHAR(20) NOT NULL,             
    email VARCHAR(40) NOT NULL UNIQUE,         
    password BYTEA NOT NULL,  
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);