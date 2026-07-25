package xdbmongodb

import (
	"time"

	"github.com/zhiyunliu/glue/xdb"
)

const (
	ScopeName = "xdb-mongodb"
)

type monitorConfig struct {
	CommandAttributeDisabled bool
	ConnName                 string
	ShowQueryLog             bool
	slowThreshold            time.Duration
	proto                    string
	logger                   xdb.Logger
}
