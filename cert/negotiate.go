package cert

import (
	"bytes"
	"github.com/vmihailenco/msgpack/v5"
	"time"
)

func GenerateExchange(privateKey, PublicKey *bytes.Buffer, password []byte, exchangeInfo *ExchangeInfo) []byte {
	exchangeInfo.PublicKey = PublicKey.Bytes()
	exchangeInfo.Timestamp = time.Now().UnixNano()

	info, err := msgpack.Marshal(*exchangeInfo)
	if err != nil {
		logger.Error("failed to dump exchange info")
		return nil
	}
	// encrypt order must be maintained !!
	signature := encrypt(password, sign(privateKey, info))
	cypherText := encrypt(password, info)
	return append(append(signature, _EncryptionSepByte), cypherText...)
}

func ExtractExchangeInfo(password, cypher []byte) *ExchangeInfo {
	split := bytes.Split(cypher, []byte{_EncryptionSepByte})
	if len(split) != 2 {
		logger.Error("malformed exchange info")
		return nil
	}

	// decrypt order must be maintained !!
	signature, err := decrypt(password, split[0])
	if err != "" {
		logger.Error("failed to decrypt signature")
		return nil
	}

	plain, err := decrypt(password, split[1])
	if err != "" {
		logger.Error("failed to decrypt exchange info")
		return nil
	}
	exchangeInfo := ExchangeInfo{}

	e := msgpack.Unmarshal(plain, &exchangeInfo)
	if e != nil {
		logger.Error(e.Error())
		return nil
	}

	if !verify(bytes.NewBuffer(exchangeInfo.PublicKey), plain, signature) {
		return nil
	}
	return &exchangeInfo
}
