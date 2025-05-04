CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(20) NOT NULL UNIQUE,
    email VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(64) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE profiles_users (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    username VARCHAR(20) REFERENCES users(username),
    avatar BYTEA,
    user_bio TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



CREATE TABLE landmark (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    category VARCHAR(50), 
    description TEXT NOT NULL,
    price INTEGER,
    photo BYTEA NOT NULL, 
    latitude FLOAT NOT NULL,
    longitude FLOAT NOT NULL
);


CREATE TABLE category (
    id SERIAL PRIMARY KEY,
    landmark_id INTEGER REFERENCES landmark(id),  
    category_name VARCHAR(50) NOT NULL
);

CREATE TABLE traffic (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) 
);

CREATE TABLE landmarks_traffic (
    id SERIAL PRIMARY KEY,
    traffic_id INTEGER REFERENCES traffic(id),      
    landmark_id INTEGER REFERENCES landmark(id)     
);
