
CREATE TABLE weather (
        id SERIAL PRIMARY KEY,
        landmark_id INT REFERENCES landmark(id) ,
        date timestamptz DEFAULT current_timestamp  ,
        temperature float ,
        description text,
        icon text,
        rain float,
        wind_speed float,
        wind_degree int,
        UNIQUE(landmark_id,date)
);