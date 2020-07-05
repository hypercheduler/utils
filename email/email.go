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

func (n *Notifier) getAddress() string {
	return fmt.Sprintf("%s:%d", n.SMTPHost, n.SMTPPort)
}

// validate smtp server connection & authorize & check ssl status
// execute on every notifier change
func (n *Notifier) Validate() NotifyError {
	tlsConfig := &tls.Config{ServerName: n.SMTPHost}
	address := n.getAddress()
	n.ViaSSL = true
	tlsConn, err := tls.Dial("tcp", address, tlsConfig)
	if err != nil {
		// tls connection failed
		n.ViaSSL = false
	}
	defer func() {
		_ = tlsConn.Close()
	}()
	client := &smtp.Client{}
	if n.ViaSSL {
		client, err = smtp.NewClient(tlsConn, n.SMTPHost)
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
		client, err = smtp.NewClient(conn, n.SMTPHost)
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
	auth := smtp.PlainAuth("", n.User, n.Password, n.SMTPHost)
	if client.Auth(auth) != nil {
		return AuthorizeError
	}
	return NoError
}

func (n *Notifier) Send(title string, text string, notifyList []string, headers map[string]string) {
	if title == "" || text == "" || len(notifyList) == 0 {
		return
	}
	address := n.getAddress()
	client := &smtp.Client{}
	if !n.ViaSSL {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		client, err = smtp.NewClient(conn, n.SMTPHost)
		if err != nil {
			logger.Error(err.Error())
			return
		}
	} else {
		conn, err := tls.Dial("tcp", address, &tls.Config{ServerName: n.SMTPHost})
		if err != nil {
			logger.Error(err.Error())
			return
		}
		client, err = smtp.NewClient(conn, n.SMTPHost)
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
	auth := smtp.PlainAuth("", n.User, n.Password, n.SMTPHost)
	if client.Auth(auth) != nil {
		logger.Error("Server Auth Failed!")
		return
	}
	if client.Mail(n.User) != nil {
		logger.Error("Mail request Failed!")
		return
	}
	for _, rcpt := range notifyList {
		if client.Rcpt(rcpt) != nil {
			logger.Error("failed to set rcpt: %s", zap.String("rcpt", rcpt))
			return
		}
	}

	headers["From"] = n.User
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
