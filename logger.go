package xdbmongodb

import (
	"context"
	"fmt"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/zhiyunliu/glue/xdb"
	"go.mongodb.org/mongo-driver/bson"
)

type mongoLogger struct {
	slowCfg *monitorConfig
}

func (l mongoLogger) Info(level int, message string, keysAndValues ...interface{}) {

}

// Error logs an error message with the given key/value pairs
func (l mongoLogger) Error(err error, message string, keysAndValues ...interface{}) {
	if l.slowCfg.logger == nil {
		return
	}
	l.slowCfg.logger.Log(context.Background(), 0, l.slowCfg.ConnName, fmt.Sprintf("[%s][%d]%s,err:%s", l.slowCfg.ConnName, 0, message, err.Error()), keysAndValues...)
}

type slowConfig struct {
	cmdCache      cmap.ConcurrentMap[int64, bson.Raw]
	ConnName      string
	ShowQueryLog  bool
	logger        xdb.Logger
	slowThreshold time.Duration
}

