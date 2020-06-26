package storage

import (
	"github.com/etcd-io/bbolt"
	"github.com/hypercheduler/utils/log"
)

var logger = log.GetLogger("storage")

type Serializer struct {
	Path string
	db   *bbolt.DB
}

type KeyValueSaver struct {
	Serializer
	GroupName string
}

type CommitSaver struct {
	KeyValueSaver
}