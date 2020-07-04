package storage

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestSerializer(t *testing.T) {
	f, _ := ioutil.TempFile("", "serializer")
	defer os.Remove(f.Name())
	s := &Serializer{}
	s.Initial(f.Name())
	if !s.Set([]byte("db1"), []byte("key"), []byte("value")) ||
		string(s.Get([]byte("db1"), []byte("key"))) != "value" {
		t.Error("failed to set kv")
	}
	s.Close()
}
