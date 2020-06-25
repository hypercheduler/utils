package cert

import (
	"testing"
)

func TestGenerateExchange(t *testing.T) {
	private, public := Generate("test")
	password := "password"
	exchange := ExchangeInfo{
		ServerId: "id-001",
	}
	extract := ExtractExchangeInfo([]byte(password), GenerateExchange(private, public, []byte(password), &exchange))
	if extract.ServerId != exchange.ServerId || string(extract.PublicKey) != string(public.Bytes()) {
		t.Error("failed to exchange info")
	}
}
