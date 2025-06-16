package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"trailblazer/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserDB struct {
	ctx      context.Context
	postgres *sqlx.DB
}

func NewUserPostgres(ctx context.Context, db *sqlx.DB) *UserDB {
	return &UserDB{ctx: ctx, postgres: db}
}

func (u *UserDB) GetUser(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT * FROM users WHERE email = $1`
	var result models.User
	err := u.postgres.GetContext(ctx, &result, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &result, nil
}

func (u *UserDB) AddUser(ctx context.Context, userData models.User) error {
	tx, err := u.postgres.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`
	_, err = u.postgres.ExecContext(ctx, query, userData.Username, userData.Email, userData.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to add user %w", err)
	}

	var userID int
	err = tx.QueryRowContext(ctx, `SELECT id FROM users WHERE email = $1`, userData.Email).Scan(&userID)
	if err != nil {
		return fmt.Errorf("failed to get user id: %w", err)
	}
	profile := models.Profile{
		UserID:    userID,
		Username:  userData.Username,
		AvatarURL: nil,
		UserBIO:   "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	profileQuery := `INSERT INTO profiles_users (user_id, username, avatar, user_bio, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.ExecContext(ctx, profileQuery, profile.UserID, profile.Username, profile.AvatarURL, profile.UserBIO, profile.CreatedAt, profile.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to add profile: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (u *UserDB) GetProfile(ctx context.Context, userID int64) (*models.Profile, error) {
	query := `SELECT username, user_bio, avatar FROM profiles_users WHERE user_id = $1`
	var profile models.Profile
	profile.UserID = int(userID)
	err := u.postgres.GetContext(ctx, &profile, query, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	return &profile, nil
}

func (u *UserDB) UpdateUserProfile(ctx context.Context, userID int, username string, avatarURL []byte, userBIO string) error {
	tx, err := u.postgres.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, `UPDATE users SET username = $1 WHERE id = $2`, username, userID)
	if err != nil {
		return fmt.Errorf("update users: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE profiles_users
		SET username = $1,
			avatar = $2,
			user_bio = $3,
			updated_at = $4
		WHERE user_id = $5
	`, username, avatarURL, userBIO, time.Now(), userID)

	if err != nil {
		return fmt.Errorf("update profiles_users: %w", err)
	}

	return nil

}
func (u *UserDB) AddReview(review models.Review) error {
	tx, err := u.postgres.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	query := `
			INSERT INTO reviews(landmark_id,user_id,rating,review)
			VALUES ($1, $2, $3, $4)
			`
	_, err = tx.Exec(query, review.LandmarkID, review.UserID, review.Rating, review.Review)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error with adding review: %w", err)
	}
	query = `
			SELECT id from reviews WHERE landmark_id=$1 AND user_id=$2 ORDER BY id DESC
			`
	row := tx.QueryRow(query, review.LandmarkID, review.UserID)
	var reviewID int64
	err = row.Scan(&reviewID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error with adding review: %w", err)
	}
	for name, image := range review.Images {
		query = `
			INSERT INTO reviews_images(review_id,photo_name,photo)
			VALUES ($1, $2,$3)
			`
		_, err = tx.Exec(query, reviewID, name, image)
		if err != nil {
			slog.Error(fmt.Sprintf("error with adding review: %w", err))
		}
	}
	tx.Commit()
	return nil
}

func (u *UserDB) GetReview(name string, onlyPhoto bool) (map[int]models.ReviewByUser, error) {
	query := `
			SELECT r.id,r.rating, r.review,ri.photo,ri.photo_name, l.name,pu.username,pu.avatar
				FROM reviews as r
				JOIN public.landmark l ON l.id = r.landmark_id
				LEFT JOIN public.reviews_images ri on r.id = ri.review_id
				JOIN public.users u on u.id = r.user_id
				JOIN public.profiles_users pu on u.id = pu.user_id
				WHERE l.images_name LIKE '%' || $1 || '%';

				`
	rows, err := u.postgres.Query(query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)

	}
	defer rows.Close()
	var reviews map[int]models.ReviewByUser = make(map[int]models.ReviewByUser)
	for rows.Next() {
		var review models.ReviewByUser
		review.Images = make(map[string][]byte)
		var id int
		var photo []byte
		var photoName sql.NullString
		err = rows.Scan(&id, &review.Rating, &review.Review.Review, &photo, &photoName, &review.LandmarkName, &review.Username, &review.Avatar)
		if err != nil {
			return nil, fmt.Errorf("failed to get reviews: %w", err)
		}
		if _, ok := reviews[id]; !ok {
			if photoName.String == "" && onlyPhoto {
				continue
			}
			reviews[id] = review
			if photoName.String != "" {
				reviews[id].Images[photoName.String] = photo
			}
		} else {
			if photoName.String != "" {
				reviews[id].Images[photoName.String] = photo
			}
		}
	}
	return reviews, nil
}
func (u *UserDB) UpdateToken(token, email string) error {
	query := `INSERT INTO
    				email_tokens (email, token, expires_at,created_at) 
				VALUES ($1, $2, CURRENT_TIMESTAMP+interval '5 hours', CURRENT_TIMESTAMP) 
				ON CONFLICT (email) DO UPDATE  SET 
					token=$3, expires_at=CURRENT_TIMESTAMP+interval '5 hours', created_at=CURRENT_TIMESTAMP`
	_, err := u.postgres.Exec(query, email, token, token)
	if err != nil {
		return fmt.Errorf("Ошибка создания или обновления токена")
	}
	return nil
}

func (u *UserDB) Delete(email string) error {
	query := "DELETE FROM users WHERE email=$1"
	_, err := u.postgres.Exec(query, email)
	if err != nil {
		return fmt.Errorf("Ошибка удаления пользователя")
	}
	return nil
}
func (u *UserDB) VerifyEmail(token string) error {
	var email string
	var expires_at time.Time
	tx, err := u.postgres.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	query := "SELECT email,expires_at from email_tokens where token=$1"
	err = tx.QueryRow(query, token).Scan(&email, &expires_at)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return fmt.Errorf("Токен не найден")
		}
		return fmt.Errorf("Ошибка проверки токена")
	}
	if time.Now().After(expires_at) {
		tx.Rollback()
		return fmt.Errorf("Время действия токена истекло")
	}
	query = "UPDATE users SET verified = true WHERE email = $1"
	_, err = tx.Exec(query, email)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Ошибка обновления статуса пользователя: %v", err)
	}
	query = "DELETE FROM email_tokens WHERE token = $1"
	_, err = tx.Exec(query, token)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Ошибка удаления токена")
	}
	tx.Commit()

	return nil
}
