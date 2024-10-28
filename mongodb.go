package xdbmongodb

import (
	"context"
	"fmt"

	"github.com/zhiyunliu/glue/config"
	contribxdb "github.com/zhiyunliu/glue/contrib/xdb"
	"github.com/zhiyunliu/glue/global"
	"github.com/zhiyunliu/glue/xdb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	mgOpts, err := s.buildMongodbOpts(cfg.Cfg)
	if err != nil {
		return
	}
	client, err := mongo.Connect(ctx, mgOpts)
	if err != nil {
		return nil, fmt.Errorf("mongo.Connect(%s):%w", connName, err)
	}

	return &mongodb{
		client: client,
		cfg:    cfg,
	}, nil
}

func (s *mongoResolver) buildMongodbOpts(cfg *xdb.Config) (opts *options.ClientOptions, err error) {
	opts = options.Client()
	opts.ApplyURI(cfg.Conn)
	if err = opts.Validate(); err != nil {
		return opts, fmt.Errorf("mongo.Validate:%w", err)
	}

	if opts.AppName != nil && len(*opts.AppName) <= 0 {
		opts.SetAppName(global.AppName)
	}
	return
}
