
CREATE TABLE IF NOT EXISTS reviews(
    id SERIAL PRIMARY KEY ,
    landmark_id INT REFERENCES landmark(id),
    user_id INT REFERENCES users(id),
    rating INT NOT NULL,
    review varchar(1500)

);
CREATE TABLE IF NOT EXISTS reviews_images(
    id SERIAL PRIMARY KEY ,
    photo_name varchar(350),
    review_id INT REFERENCES reviews(id),
    photo BYTEA
)


