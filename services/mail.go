package services

import (
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"crypto/tls"
)

type Mailer interface {
	SendMail(fromAddress string, toAddress string, subject string, body string, smtpHostPort string, username string, password string)
}

type MailerImpl struct{}

// SSL/TLS Email Example
// https://gist.github.com/chrisgillis/10888032
func (*MailerImpl) SendMail(fromAddress string, toAddress string, subject string, body string, smtpHostPort string, username string, password string) {

	from := mail.Address{"", fromAddress}
	to   := mail.Address{"", toAddress}
	subj := subject

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k,v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	host, _, _ := net.SplitHostPort(smtpHostPort)

	auth := smtp.PlainAuth("",username, password, host)

	// TLS config
	tlsconfig := &tls.Config {
		InsecureSkipVerify: true,
		ServerName: host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", smtpHostPort, tlsconfig)
	defer conn.Close()
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	defer c.Close()
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := c.Data()
	defer w.Close()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	c.Quit()

}