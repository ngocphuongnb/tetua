package mail

import (
	"crypto/tls"
	"fmt"
	"mime"
	"net/mail"
	"net/smtp"

	"github.com/ngocphuongnb/tetua/app/config"
)

func Send(receiverName, receiverAddress, subject, body string) error {
	if config.Mail == nil ||
		config.Mail.Sender == "" ||
		config.Mail.Smtp == nil ||
		config.Mail.Smtp.Host == "" ||
		config.Mail.Smtp.Port == 0 ||
		config.Mail.Smtp.User == "" ||
		config.Mail.Smtp.Pass == "" {
		return fmt.Errorf("Mail config is not set")
	}

	from := mail.Address{
		Name:    config.Setting("app_name"),
		Address: config.Mail.Sender,
	}
	to := mail.Address{
		Name:    receiverName,
		Address: receiverAddress,
	}
	tlsconfig := &tls.Config{
		InsecureSkipVerify: config.DEVELOPMENT,
		ServerName:         config.Mail.Smtp.Host,
	}

	message := fmt.Sprintf("%s: %s\r\n", "From", from.String())
	message += fmt.Sprintf("%s: %s\r\n", "To", to.String())
	message += fmt.Sprintf("%s: %s\r\n", "Subject", mime.QEncoding.Encode("UTF-8", subject))
	message += "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	message += "\r\n" + body
	auth := smtp.PlainAuth("", config.Mail.Smtp.User, config.Mail.Smtp.Pass, config.Mail.Smtp.Host)
	address := fmt.Sprintf("%s:%d", config.Mail.Smtp.Host, config.Mail.Smtp.Port)

	tlsConn, err := tls.Dial("tcp", address, tlsconfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(tlsConn, config.Mail.Smtp.Host)
	if err != nil {
		return err
	}

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(from.Address); err != nil {
		return err
	}

	if err = client.Rcpt(to.Address); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	if _, err = w.Write([]byte(message)); err != nil {
		return err
	}

	if err = w.Close(); err != nil {
		return err
	}

	client.Quit()

	return nil
}
