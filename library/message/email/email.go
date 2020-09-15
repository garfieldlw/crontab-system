package email

import (
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"strings"
)

const (
	user     = "username@email.com"
	password = "password"
	host     = "smtp.qiye.aliyun.com:465"
)

func SendMail(to, subject, body, mailType string) error {
	contentType := ""
	if mailType == "html" {
		contentType = "Content-Type: text/" + mailType + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	sendTo := strings.Split(to, ";")

	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	err := smtp.SendMail(host, auth, user, sendTo, msg)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func SendMailTLS(to, subject, body, mailType string) error {
	sendTo := strings.Split(to, ";")
	e := email.NewEmail()
	e.From = user
	e.To = sendTo
	e.Subject = subject
	if mailType == "html" {
		e.HTML = []byte(body[:])
	} else {
		e.Text = []byte(body[:])
	}

	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	return e.SendWithTLS(host, auth, &tls.Config{ServerName: hp[0]})
}
