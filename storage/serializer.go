package storage

import (
	"github.com/etcd-io/bbolt"
	"go.uber.org/zap"
	"os"
	"time"
)

func (self *Serializer) Initial(path string, readOnly bool) (ok bool) {
	options := &bbolt.Options{Timeout: time.Second}

	var mode os.FileMode
	if readOnly {
		mode = 0400
		options.ReadOnly = true
	} else {
		mode = 0600
	}

	self.Path = path
	db, err := bbolt.Open(self.Path, mode, options)
	if err != nil {
		logger.Error(err.Error())
		return false
	}
	self.db = db
	return true
}

func (self *Serializer) Save(group, key, value string, initBucket bool) bool {
	if self.db.Batch(func(tx *bbolt.Tx) error {
		var bucket = &bbolt.Bucket{}
		if initBucket {
			bucket, _ = tx.CreateBucket([]byte(group))
		} else {
			bucket = tx.Bucket([]byte(group))
		}
		return bucket.Put([]byte(key), []byte(value))
	}) != nil {
		logger.Error("failed to update",
			zap.String("group", group),
			zap.String("key", key))
		return false
	}
	return true
}

func (self *Serializer) Read(group, key string) (value string) {
	if self.db.View(func(tx *bbolt.Tx) error {
		value = string(tx.Bucket([]byte(group)).Get([]byte(key)))
		return nil
	}) != nil {
		logger.Error("failed to read",
			zap.String("group", group),
			zap.String("key", key))
		return ""
	}
	return
}
