package email

import (
	"crypto/tls"
	"fmt"
	"go.uber.org/zap"
	"net"
	"net/smtp"
	"os"
	"strings"
)

func (self *Notifier) getAddress() string {
	return fmt.Sprintf("%s:%d", self.SMTPHost, self.SMTPPort)
}

// validate smtp server connection & authorize & check ssl status
// execute on every notifier change
func (self *Notifier) Validate() NotifyError {
	tlsConfig := &tls.Config{ServerName: self.SMTPHost}
	address := self.getAddress()
	self.ViaSSL = true
	tlsConn, err := tls.Dial("tcp", address, tlsConfig)
	if err != nil {
		// tls connection failed
		self.ViaSSL = false
	}
	defer func() {
		_ = tlsConn.Close()
	}()
	client := &smtp.Client{}
	if self.ViaSSL {
		client, err = smtp.NewClient(tlsConn, self.SMTPHost)
		if err != nil {
			return ProtocolError
		}
	} else {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			return ConnectionError
		}
		defer func() {
			_ = conn.Close()
		}()
		client, err = smtp.NewClient(conn, self.SMTPHost)
		if err != nil {
			return ProtocolError
		}
	}
	defer func() {
		_ = client.Quit()
	}()
	hostname, _ := os.Hostname()
	if client.Hello(hostname) != nil {
		return ProtocolError
	}
	auth := smtp.PlainAuth("", self.User, self.Password, self.SMTPHost)
	if client.Auth(auth) != nil {
		return AuthorizeError
	}
	return NoError
}

func (self *Notifier) Send(title string, text string, notifyList []string, headers map[string]string) {
	if title == "" || text == "" || len(notifyList) == 0 {
		return
	}
	address := self.getAddress()
	client := &smtp.Client{}
	if !self.ViaSSL {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		client, err = smtp.NewClient(conn, self.SMTPHost)
		if err != nil {
			logger.Error(err.Error())
			return
		}
	} else {
		conn, err := tls.Dial("tcp", address, &tls.Config{ServerName: self.SMTPHost})
		if err != nil {
			logger.Error(err.Error())
			return
		}
		client, err = smtp.NewClient(conn, self.SMTPHost)
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}
	defer func() {
		_ = client.Quit()
	}()
	hostname, _ := os.Hostname()
	if client.Hello(hostname) != nil {
		logger.Error("Server Hello Failed!")
		return
	}
	auth := smtp.PlainAuth("", self.User, self.Password, self.SMTPHost)
	if client.Auth(auth) != nil {
		logger.Error("Server Auth Failed!")
		return
	}
	if client.Mail(self.User) != nil {
		logger.Error("Mail request Failed!")
		return
	}
	for _, rcpt := range notifyList {
		if client.Rcpt(rcpt) != nil {
			logger.Error("failed to set rcpt: %s", zap.String("rcpt", rcpt))
			return
		}
	}

	headers["From"] = self.User
	headers["Subject"] = title
	entity := "To: " + strings.Join(notifyList, ",") + "\r\n"

	for header, value := range headers {
		entity += fmt.Sprintf("%s: %s\r\n", header, value)
	}
	entity += "\r\n" + text

	writer, err := client.Data()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	_, _ = writer.Write([]byte(entity))

	err = writer.Close()
	if err != nil {
		logger.Error(err.Error())
		return
	}

}
