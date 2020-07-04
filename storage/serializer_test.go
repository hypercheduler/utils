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

func TestSearch(t *testing.T) {
	f, _ := ioutil.TempFile("", "serializer-search")
	defer os.Remove(f.Name())
	s := &Serializer{}
	s.Initial(f.Name())
	gName := []byte("group")
	keyList := map[string]string{"h-1": "1", "h1": "2", "h0": "3", "well": "4"}
	for k, v := range keyList {
		s.Set(gName, []byte(k), []byte(v))
	}
	for _, v := range s.SearchKey(gName, []byte("h")) {
		if keyList[string(v)] == "" {
			t.Error("search error")
		}
	}
	for k, v := range s.SearchWithContent(gName, []byte("h")) {
		if keyList[k] != string(v) {
			t.Error("search content not match")
		}
	}
	if s.SearchKey(gName, []byte("H")) != nil {
		t.Error("should be nothing to match")
	}
}
