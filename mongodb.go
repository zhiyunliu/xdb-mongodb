package xdbmongodb

import (
	"context"
	"fmt"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/zhiyunliu/glue/config"
	contribxdb "github.com/zhiyunliu/glue/contrib/xdb"
	"github.com/zhiyunliu/glue/global"
	"github.com/zhiyunliu/glue/xdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

func init() {
	xdb.Register(&mongoResolver{})
	//tpl.Register(New(Proto, ArgumentPrefix))
}

const (
	Proto          = "mongodb"
	ArgumentPrefix = ""
)

type mongoResolver struct {
}

func (s *mongoResolver) Name() string {
	return Proto
}

func (s *mongoResolver) Resolve(connName string, setting config.Config, opts ...xdb.Option) (dbObj interface{}, err error) {
	cfg := contribxdb.NewConfig(connName)
	err = setting.ScanTo(cfg.Cfg)
	if err != nil {
		return nil, fmt.Errorf("读取DB配置(%s):%w", connName, err)
	}

	newCfg, err := xdb.DefaultRefactor(cfg.ConnName, cfg.Cfg)
	if err != nil {
		return
	}
	if newCfg != nil {
		cfg.Cfg = newCfg
	}

	for i := range opts {
		opts[i](cfg.Cfg)
	}

	ctx := context.Background()

	mgOpts, databaseName, err := s.buildMongodbOpts(connName, cfg.Cfg)
	if err != nil {
		return
	}
	client, err := mongo.Connect(ctx, mgOpts)
	if err != nil {
		return nil, fmt.Errorf("mongo.Connect(%s):%w", connName, err)
	}

	return &mongodb{
		client:   client,
		database: client.Database(databaseName),
		cfg:      cfg,
	}, nil
}

func (s *mongoResolver) buildMongodbOpts(connName string, cfg *xdb.Config) (opts *options.ClientOptions, databaseName string, err error) {
	cs, err := connstring.ParseAndValidate(cfg.Conn)
	if err != nil {
		err = fmt.Errorf("connstring.Validate:%w", err)
		return
	}
	databaseName = cs.Database
	opts = options.Client()
	opts.ApplyURI(cfg.Conn)
	if err = opts.Validate(); err != nil {
		return opts, databaseName, fmt.Errorf("mongo.Validate:%w", err)
	}
	if opts.AppName != nil && len(*opts.AppName) <= 0 {
		opts.SetAppName(global.AppName)
	}

	shardCnt := int64(cmap.SHARD_COUNT)

	slowCfg := &slowConfig{
		ConnName:      connName,
		ShowQueryLog:  cfg.ShowQueryLog,
		slowThreshold: time.Duration(cfg.LongQueryTime) * time.Millisecond,
		cmdCache: cmap.NewWithCustomShardingFunction[int64, bson.Raw](func(key int64) uint32 {
			return uint32(key % shardCnt)
		}),
	}
	slowCfg.logger, _ = xdb.GetLogger(cfg.LoggerName)

	opts.Monitor = &event.CommandMonitor{
		Started: func(ctx context.Context, cse *event.CommandStartedEvent) {
			slowCfg.Set(cse.RequestID, cse.Command)
		},
		Succeeded: func(ctx context.Context, cse *event.CommandSucceededEvent) {
			slowCfg.printSlowQuery(ctx, cse.RequestID, cse.Duration, fmt.Sprintf("%s.%s", cse.DatabaseName, cse.CommandName))
		},
		Failed: func(ctx context.Context, cse *event.CommandFailedEvent) {
			slowCfg.printSlowQuery(ctx, cse.RequestID, cse.Duration, fmt.Sprintf("%s.%s", cse.DatabaseName, cse.CommandName))
		},
	}
	opts.SetLoggerOptions(&options.LoggerOptions{
		Sink: &mongoLogger{},
		ComponentLevels: map[options.LogComponent]options.LogLevel{
			options.LogComponentAll: options.LogLevelDebug,
		},
	})
	return
}
