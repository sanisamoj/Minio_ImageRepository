package main

import (
	"crypto/tls"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

type InitialMailConfig struct {
	Email    string
	Password string
	Host     string
	Port     int
}

var mailConfig InitialMailConfig

func init() {
	portStr := os.Getenv("EMAIL_PORT")
	port, _ := strconv.Atoi(portStr)

	mailConfig = InitialMailConfig{
		Email:    os.Getenv("EMAIL_AUTH_USER"),
		Password: os.Getenv("EMAIL_AUTH_PASS"),
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     port,
	}
}

func sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mailConfig.Email)
	m.SetHeader("To", to)
	m.SetAddressHeader("Cc", mailConfig.Email, "Sanisamoj")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(mailConfig.Host, mailConfig.Port, mailConfig.Email, mailConfig.Password)
	d.TLSConfig = &tls.Config{ServerName: mailConfig.Host}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendLoginCodeEmail(to, username, userEmail, code string, expirationMinutes int) error {
	sub := "Seu código de login"
	body, err := convLoginCodeEmail(username, userEmail, code, expirationMinutes)
	if err != nil {
		return err
	}

	return sendEmail(to, sub, body)
}

func convLoginCodeEmail(username, userEmail, code string, expirationMinutes int) (string, error) {
	templateBytes, err := os.ReadFile("login_code.html")
	if err != nil {
		return "", err
	}

	htmlTemplate := string(templateBytes)
	body := htmlTemplate
	body = strings.ReplaceAll(body, "{PROJECT_NAME}", "Storage-repo")
	body = strings.ReplaceAll(body, "{USER_NAME}", username)
	body = strings.ReplaceAll(body, "{USER_EMAIL}", userEmail)
	body = strings.ReplaceAll(body, "{LOGIN_CODE}", code)
	body = strings.ReplaceAll(body, "{EXPIRATION_MINUTES}", strconv.Itoa(expirationMinutes))


	actYear := time.Now().Year()
	body = strings.ReplaceAll(body, "{CURRENT_YEAR}", strconv.Itoa(actYear))

	return body, nil
}
