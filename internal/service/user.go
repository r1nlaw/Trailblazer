package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log/slog"
	"math/rand"
	"strconv"

	"trailblazer/internal/config"
	"trailblazer/internal/models"
	"trailblazer/internal/repository"

	"github.com/k3a/html2text"

	"gopkg.in/gomail.v2"

	"github.com/google/uuid"
)

type User struct {
	repo       repository.User
	SMTPconfig config.SMTPConfig
	HostConfig config.HostConfig
}

func NewUserService(repository repository.User, smtpConfig config.SMTPConfig, hostConfig config.HostConfig) *User {
	return &User{
		repo:       repository,
		SMTPconfig: smtpConfig,
		HostConfig: hostConfig,
	}
}
func (s *User) GetUser(ctx context.Context, email string) (*models.User, error) {
	u, err := s.repo.GetUser(ctx, email)
	if err != nil {
		return nil, err
	}
	if u.Verified != true {
		return nil, fmt.Errorf("user is not verified")
	}
	return u, nil
}

func (s *User) AddUser(c context.Context, user models.User) error {
	return s.repo.AddUser(c, user)
}

func (s *User) AddReview(review models.Review) error {
	return s.repo.AddReview(review)
}

func (s *User) GetReview(name string, onlyPhoto bool) (map[int]models.ReviewByUser, error) {
	return s.repo.GetReview(name, onlyPhoto)
}
func (s *User) UpdateUserProfile(c context.Context, i int, username string, bytes []byte, bio string) error {
	return s.repo.UpdateUserProfile(c, i, username, bytes, bio)
}

func (s *User) VerifyEmail(token string) error {
	return s.repo.VerifyEmail(token)

}

func (s *User) SendToken(email string) error {
	token := getToken()
	subject := "Confirm Your Email"
	htmlBody := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<style>
			body {
				font-family: Arial, sans-serif;
				line-height: 1.6;
				color: #333333;
				margin: 0;
				padding: 0;
			}
			.container {
				max-width: 600px;
				margin: 0 auto;
				padding: 20px;
			}
			.header {
				background-color: #4CAF50;
				color: white;
				padding: 20px;
				text-align: center;
				border-radius: 5px 5px 0 0;
			}
			.content {
				background-color: #ffffff;
				padding: 20px;
				border: 1px solid #dddddd;
				border-radius: 0 0 5px 5px;
			}
			.button {
				display: inline-block;
				background-color: #4CAF50;
				color: white;
				padding: 12px 24px;
				text-decoration: none;
				border-radius: 4px;
				margin: 20px 0;
			}
			.footer {
				text-align: center;
				margin-top: 20px;
				color: #666666;
				font-size: 12px;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>Подтверждение Email</h1>
			</div>
			<div class="content">
				<p>Здравствуйте!</p>
				<p>Спасибо за регистрацию. Для завершения процесса регистрации, пожалуйста, подтвердите ваш email адрес.</p>
				<p>Нажмите на кнопку ниже для подтверждения:</p>
				<div style="text-align: center;">
					<a href="` + fmt.Sprintf("%s/user/verify?token=%s", s.HostConfig.Domain, token) + `" 
					   style="display: inline-block; 
					   background-color: #4CAF50; 
					   color: white; 
					   padding: 12px 24px; 
					   text-decoration: none; 
					   border-radius: 4px; 
					   font-weight: bold;
					   margin: 20px 0;
					   box-shadow: 0 2px 4px rgba(0,0,0,0.2);">
					   Подтвердить Email
					</a>
				</div>
				<p>Если вы не запрашивали это письмо, пожалуйста, проигнорируйте его.</p>
				<p>С уважением,<br>Команда Trailblazer</p>
			</div>
			<div class="footer">
				<p>Это автоматическое сообщение, пожалуйста, не отвечайте на него.</p>
			</div>
		</div>
	</body>
	</html>`
	m := gomail.NewMessage()
	m.SetHeader("From", s.SMTPconfig.Email)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)
	m.AddAlternative("text/plain", html2text.HTML2Text(htmlBody))

	port, err := strconv.Atoi(s.SMTPconfig.Port)
	if err != nil {

		return err
	}
	certPool := x509.NewCertPool()
	pemData, err := ioutil.ReadFile("smtp-cert.pem") // Убедитесь, что файл доступен
	if err != nil {
		slog.Warn("failed to read certificate file, using insecure mode", "error", err)
	} else if !certPool.AppendCertsFromPEM(pemData) {
		slog.Warn("failed to append certificate to pool")
	}
	tlsConfig := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}
	dialer := gomail.NewDialer(s.SMTPconfig.Host, port, s.SMTPconfig.Username, s.SMTPconfig.Password)
	dialer.TLSConfig = tlsConfig

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}
	err = s.repo.UpdateToken(token, email)
	if err != nil {
		return err
	}
	return nil
}

func (s *User) Delete(email string) error {
	return s.repo.Delete(email)
}

func (s *User) GetProfile(c context.Context, userID int64) (*models.Profile, error) {
	return s.repo.GetProfile(c, userID)
}
func getToken() string {
	uuid.SetRand(rand.New(rand.NewSource(1)))
	return uuid.New().String()
}
