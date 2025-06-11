CREATE TABLE landmark (
                          id SERIAL PRIMARY KEY ,
                          name VARCHAR(150) NOT NULL,
                          address text,
                          category VARCHAR(50),
                          description TEXT,
                          history TEXT ,
                          photo BYTEA ,
                          location geography(POINT, 4326)

);



