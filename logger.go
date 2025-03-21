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
	slowCfg *slowConfig
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

func (slowCfg *slowConfig) printSlowQuery(ctx context.Context, requestId int64, timeRange time.Duration, query string) {
	if !slowCfg.ShowQueryLog {
		return
	}
	if slowCfg.logger == nil {
		return
	}

	queryRaw, ok := slowCfg.cmdCache.Get(requestId)
	if !ok {
		return
	}
	slowCfg.cmdCache.Remove(requestId)

	if slowCfg.slowThreshold <= 0 || timeRange < slowCfg.slowThreshold {
		return
	}
	slowCfg.logger.Log(ctx, timeRange.Milliseconds(), slowCfg.ConnName, fmt.Sprintf("[%s][%d]%s", slowCfg.ConnName, requestId, query), queryRaw)
}

func (slowCfg *slowConfig) Set(requestId int64, cmd bson.Raw) {
	if !slowCfg.ShowQueryLog {
		return
	}
	if slowCfg.logger == nil {
		return
	}

	slowCfg.cmdCache.Set(requestId, cmd)
}
