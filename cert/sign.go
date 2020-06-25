package cert

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"math/big"
)

func createPrivateKey(privateKey *bytes.Buffer) *ecdsa.PrivateKey {
	block, _ := pem.Decode(privateKey.Bytes())
	if block == nil {
		logger.Error("private key format error")
		return nil
	}
	pvKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		logger.Error("failed to parse private key")
	}
	return pvKey
}

func createPublicKey(publicKey *bytes.Buffer) *ecdsa.PublicKey {
	block, _ := pem.Decode(publicKey.Bytes())
	if block == nil {
		logger.Error("public key format error")
		return nil
	}
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		logger.Error("failed to parse public key")
		return nil
	}
	pk, ok := certificate.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		logger.Error("this is not a ecdsa public key")
		return nil
	}
	return pk
}

func sign(privateKey *bytes.Buffer, src []byte) []byte {
	r, s, err := ecdsa.Sign(rand.Reader, createPrivateKey(privateKey), src)
	if err != nil {
		logger.Error("failed to sign src")
		return nil
	}
	rText, _ := r.MarshalText()
	sText, _ := s.MarshalText()
	return append(append(rText, _SignatureSepByte), sText...)
}

func verify(publicKey *bytes.Buffer, src, signature []byte) bool {
	var rInt, sInt big.Int
	split := bytes.Split(signature, []byte{_SignatureSepByte})
	if len(split) != 2 {
		return false
	}
	eR := rInt.UnmarshalText(split[0])
	eS := sInt.UnmarshalText(split[1])
	if eR != nil || eS != nil {
		logger.Error("failed to parse signature")
		return false
	}
	return ecdsa.Verify(createPublicKey(publicKey), src, &rInt, &sInt)
}

func encrypt(password, src []byte) []byte {
	var passwordLen = len(password)

	var extraTable = []byte(_EncryptionTable)
	var extraTableLen = len(extraTable)

	for index, bt := range password {
		// warning !! outer password will be changed
		password[index] = bt + extraTable[index%extraTableLen]
	}

	for index, bt := range src {
		src[index] = bt ^ password[index%passwordLen] + extraTable[index%extraTableLen]
	}
	return []byte(base64.StdEncoding.EncodeToString(src))
}

func decrypt(password, cypher []byte) ([]byte, string) {
	tmp, err := base64.StdEncoding.DecodeString(string(cypher))
	if err != nil {
		return nil, "failed to decode base64"
	}

	var extraTable = []byte(_EncryptionTable)
	var extraTableLen = len(extraTable)
	var passwordLen = len(password)

	for index, bt := range password {
		// warning !! outer password will be changed
		password[index] = bt + extraTable[index%extraTableLen]
	}

	for index, bt := range tmp {
		tmp[index] = (bt - extraTable[index%extraTableLen]) ^ password[index%passwordLen]
	}
	return tmp, ""
}
