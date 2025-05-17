CREATE TABLE landmark (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(150) NOT NULL,
                          address text,
                          category VARCHAR(50),
                          description TEXT,
                          history TEXT ,
                          photo BYTEA ,
                          location geography(POINT, 4326)

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

