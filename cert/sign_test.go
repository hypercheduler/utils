package cert

import "testing"

func TestSign(t *testing.T) {
	originText := []byte("hellflame is fine")
	privateKey, publicKey := Generate("localhost")
	signature := sign(privateKey, originText)
	if verify(publicKey, []byte("well"), signature) {
		t.Error("they actually don't match")
	}
	if !verify(publicKey, originText, signature) {
		t.Error("they actually match")
	}
}

func TestEncrypt(t *testing.T) {
	testGroup := map[string]string{
		"%4":    "just test",
		"09123": "901821",
		"":      "",
		"aaaa":  "bbbb",
		"bbbb":  "aaaa",
		"cccc":  "cccc",
		"dddd":  "dddd",
		"ab":    "ba",
		"ba":    "ab",
	}

	for password, plain := range testGroup {
		extract, err := decrypt([]byte(password), encrypt([]byte(password), []byte(plain)))
		if err != "" {
			t.Error(err)
		}
		if string(extract) != plain {
			t.Errorf("%s != %s", extract, plain)
		}
	}
}
