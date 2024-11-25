package main

import (
	"context"

	"github.com/zhiyunliu/glue"
	_ "github.com/zhiyunliu/glue/contrib/metrics/prometheus"
	"github.com/zhiyunliu/glue/global"
	"github.com/zhiyunliu/glue/log"
	"github.com/zhiyunliu/glue/server/api"
	"github.com/zhiyunliu/glue/xdb"
	_ "github.com/zhiyunliu/xdb-mongodb"
	"github.com/zhiyunliu/xdb-mongodb/samples/services"
)

func main() {

	global.AppName = "mongodb-sample"
	apiSrv := api.New("apiserver", api.WithServiceName(global.AppName), api.Log(log.WithRequest(), log.WithResponse()))
	services.BindAPI(apiSrv)

	opts := []glue.Option{glue.Server(apiSrv), glue.StartingHook(func(ctx context.Context) error {

		xdb.Default.ShowQueryLog = true
		xdb.Default.LongQueryTime = 1
		xdb.RegistryLogger(&dbLogger{
			name: "dbslowsql",
		})
		//return global.Config.ScanTo(config.Sys)
		return nil
	})}

	app := glue.NewApp(opts...)
	app.Start()
}
