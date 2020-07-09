package email

import (
	"github.com/hypercheduler/utils"
	"github.com/hypercheduler/utils/log"
)

type Notifier struct {
	User     string
	Password string
	SMTPHost string
	SMTPPort int

	ViaSSL bool
}

var logger = log.GetLogger("utils-email", utils.VERSION)

type NotifyError int

const (
	NoError NotifyError = iota
	ConnectionError
	ProtocolError
	AuthorizeError
)
