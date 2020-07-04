package storage

import (
	"github.com/etcd-io/bbolt"
	"github.com/hypercheduler/utils"
	"github.com/hypercheduler/utils/log"
)

var logger = log.GetLogger("storage", utils.VERSION)

type Serializer struct {
	Path    string
	HasInit bool
	db      *bbolt.DB
}
