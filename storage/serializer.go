package storage

import (
	"bytes"
	"github.com/etcd-io/bbolt"
	"go.uber.org/zap"
	"time"
)

func (s *Serializer) Initial(path string) (ok bool) {
	s.Path = path

	db, err := bbolt.Open(s.Path, 0600,
		&bbolt.Options{Timeout: time.Second})
	if err != nil {
		logger.Error("failed to open db",
			zap.String("error", err.Error()),
			zap.String("file", s.Path))
		return false
	}
	s.db = db
	return true
}

func (s *Serializer) Set(group, key, value []byte) bool {
	if s.db.Batch(func(tx *bbolt.Tx) error {
		bucket, _ := tx.CreateBucketIfNotExists(group)
		return bucket.Put(key, value)
	}) != nil {
		logger.Error("failed to update",
			zap.String("group", string(group)),
			zap.String("key", string(key)))
		return false
	}
	return true
}

func (s *Serializer) Get(group, key []byte) (value []byte) {
	if s.db.View(func(tx *bbolt.Tx) error {
		value = tx.Bucket(group).Get(key)
		return nil
	}) != nil {
		logger.Error("failed to read",
			zap.String("group", string(group)),
			zap.String("key", string(key)))
		return nil
	}
	return
}

func (s *Serializer) SearchKey(group, prefix []byte) (keyList [][]byte) {
	if s.db.View(func(tx *bbolt.Tx) error {
		cursor := tx.Bucket(group).Cursor()
		for k, _ := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = cursor.Next() {
			keyList = append(keyList, k)
		}
		return nil
	}) != nil {
		logger.Error("failed to search prefix",
			zap.String("group", string(group)),
			zap.String("prefix", string(prefix)))
	}
	return
}

func (s *Serializer) SearchWithContent(group, prefix []byte) map[string][]byte {
	var result = make(map[string][]byte)
	if s.db.View(func(tx *bbolt.Tx) error {
		cursor := tx.Bucket(group).Cursor()
		for k, v := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = cursor.Next() {
			result[string(k)] = v
		}
		return nil
	}) != nil {
		logger.Error("failed to search with content",
			zap.String("group", string(group)),
			zap.String("prefix", string(prefix)))
	}
	return result
}

func (s *Serializer) Close() {
	s.db.Close()
}
