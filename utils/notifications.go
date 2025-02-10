package utils

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"time"

	"github.com/FiodhyAN/learn-rest-api-golang/config"
	"github.com/FiodhyAN/learn-rest-api-golang/helpers"
	"github.com/FiodhyAN/learn-rest-api-golang/types"
	"github.com/google/uuid"
)

func SendVerificationMail(store types.UserStore, user *types.User) error {
	encryptedID, err := EncryptText(user.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	uuidV7, err := uuid.NewV7()
	if err != nil {
		return err
	}
	verificationToken := uuidV7.String() + user.ID
	encryptedVerificationToken, err := EncryptText(verificationToken)
	if err != nil {
		return err
	}
	verificationLink := config.Envs.FrontendUrl + `/verify-email?userId=` + encryptedID + `&token=` + encryptedVerificationToken

	expirationDate := time.Now().AddDate(0, 0, 1)
	finalFormattedDate, err := FormatDate(expirationDate)
	if err != nil {
		return err
	}

	templateData := struct {
		VerificationLink string
		ExpirationDate   string
	}{
		VerificationLink: verificationLink,
		ExpirationDate:   finalFormattedDate,
	}

	toList := []string{user.Email}

	headers := make(map[string]string)
	headers["From"] = config.Envs.SMTPEmail
	headers["To"] = user.Email
	headers["Subject"] = "Email Verification!"
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	var headerString string
	for key, value := range headers {
		headerString += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	t, err := template.ParseFiles("template/email-verification.html")
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, templateData); err != nil {
		return err
	}

	message := headerString + "\r\n" + buf.String()
	msg := []byte(message)

	mail_host := config.Envs.SMTPHost
	mail_port := config.Envs.SMTPPort
	mail_username := config.Envs.SMTPUsername
	mail_from := config.Envs.SMTPEmail
	mail_password := config.Envs.SMTPPassword

	mail_auth := smtp.PlainAuth("", mail_username, mail_password, mail_host)

	if err := smtp.SendMail(mail_host+":"+mail_port, mail_auth, mail_from, toList, msg); err != nil {
		return err
	}

	hashedToken, err := helpers.CreateHashPassword(verificationToken)
	if err != nil {
		return err
	}

	ctx := context.Background()
	userUUID, err := uuid.Parse(user.ID)
	if err != nil {
		return err
	}
	if err := store.UpdateUserVerificationExpired(ctx, userUUID, expirationDate, hashedToken); err != nil {
		return err
	}

	return nil
}
