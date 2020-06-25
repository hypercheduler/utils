package storage

import (
	"github.com/etcd-io/bbolt"
)

func (self *KeyValueSaver) Initial(path, groupName string, readOnly bool) (ok bool) {
	ok = self.Serializer.Initial(path, readOnly)
	if !ok {
		return false
	}
	self.GroupName = groupName
	if self.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucket([]byte(groupName))
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		return nil
	}) != nil {
		return false
	}
	return true
}

func (self *KeyValueSaver) Save(key, value string) bool {
	return self.Serializer.Save(self.GroupName, key, value, false)
}

func (self *KeyValueSaver) Read(key string) string {
	return self.Serializer.Read(self.GroupName, key)
}
