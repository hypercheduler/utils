package cert

import (
	"github.com/hypercheduler/utils"
	"github.com/hypercheduler/utils/log"
	"os"
)

var logger = log.GetLogger("utils-cert", utils.VERSION)

const _SignatureSepByte = byte(3)
const _EncryptionSepByte = byte(10)

var _EncryptionTable = "\xb2wtS,\xce\x11\xa9\xfc\xc1\xff?j\xb4\x84@"

func init() {
	fromEnv := os.Getenv("UTILS_ENC_TABLE")
	if fromEnv != "" {
		_EncryptionTable = fromEnv
	}
}

type ExchangeInfo struct {
	PublicKey   []byte
	ServerId    string
	Timestamp   int64
	UtilVersion string
}
