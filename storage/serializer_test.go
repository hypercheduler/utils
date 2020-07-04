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

	tests := []struct {
		HasContent bool
		Shred      string
		FullSearch bool
		Length     int
	}{{false, "h", false, 3},
		{false, "h", true, 3},
		{false, "e", true, 1},
		{false, "e", false, 0},
		{false, "1", true, 2},
		{true, "h", false, 3}}
	for _, v := range tests {
		var r int
		if v.HasContent {
			r = len(s.SearchWithContent(gName, []byte(v.Shred), v.FullSearch))
		} else {
			r = len(s.SearchKey(gName, []byte(v.Shred), v.FullSearch))
		}
		if r != v.Length {
			t.Error("search count not match")
		}
	}
	for _, v := range s.SearchKey(gName, []byte("h"), false) {
		if keyList[string(v)] == "" {
			t.Error("search error")
		}
	}
	for k, v := range s.SearchWithContent(gName, []byte("h"), false) {
		if keyList[k] != string(v) {
			t.Error("search content not match")
		}
	}
	if s.SearchKey(gName, []byte("H"), false) != nil {
		t.Error("should be nothing to match")
	}

}
