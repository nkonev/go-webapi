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
	SendMail(toAddress string, subject string, body string)
}

type MailerImpl struct{
	FromAddress, SmtpHostPort, Username, Password string
}

// SSL/TLS Email Example
// https://gist.github.com/chrisgillis/10888032
func (m *MailerImpl) SendMail(toAddress string, subject string, body string) {

	from := mail.Address{"", m.FromAddress}
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
	host, _, _ := net.SplitHostPort(m.SmtpHostPort)

	auth := smtp.PlainAuth("", m.Username, m.Password, host)

	// TLS config
	tlsconfig := &tls.Config {
		InsecureSkipVerify: true,
		ServerName: host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", m.SmtpHostPort, tlsconfig)
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