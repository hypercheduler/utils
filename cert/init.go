package cert

import (
	"github.com/hypercheduler/utils"
	"github.com/hypercheduler/utils/log"
)

var logger = log.GetLogger("cert", utils.VERSION)

const _SignatureSepByte = byte(3)
const _EncryptionSepByte = byte(10)
const _EncryptionTable = "\xb2wtS,\xce\x11\xa9\xfc\xc1\xff?j\xb4\x84@"

type ExchangeInfo struct {
	PublicKey   []byte
	ServerId    string
	Timestamp   int64
	UtilVersion string
}
