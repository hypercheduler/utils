package email

import (
	"github.com/hypercheduler/utils"
	"github.com/hypercheduler/utils/log"
)

type Notifier struct {
	User     string `json:"user"`
	Password string `json:"password"`
	SMTPHost string `json:"host"`
	SMTPPort int    `json:"port"`

	ViaSSL bool
}

var logger = log.GetLogger("email", utils.VERSION)

type NotifyError int

const (
	NoError NotifyError = iota
	ConnectionError
	ProtocolError
	AuthorizeError
)
