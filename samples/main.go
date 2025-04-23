package main

import (
	sctx "context"

	"github.com/zhiyunliu/glue"
	_ "github.com/zhiyunliu/glue/contrib/metrics/prometheus"
	"github.com/zhiyunliu/glue/global"
	"github.com/zhiyunliu/glue/log"
	"github.com/zhiyunliu/glue/opentelemetry"
	"github.com/zhiyunliu/glue/server/api"
	"github.com/zhiyunliu/glue/xdb"
	_ "github.com/zhiyunliu/xdb-mongodb"
	"github.com/zhiyunliu/xdb-mongodb/samples/services"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	provider *sdktrace.TracerProvider
)

func main() {

	global.AppName = "mongodb-sample"
	apiSrv := api.New("apiserver", api.WithServiceName(global.AppName), api.Log(log.WithRequest(), log.WithResponse()))
	services.BindAPI(apiSrv)

	app := glue.NewApp(glue.Server(apiSrv), glue.StartingHook(StartingHook), glue.StopedHook(StopedHook))
	app.Start()
}

func StartingHook(ctx sctx.Context) (err error) {

	xdb.Default.ShowQueryLog = true
	xdb.Default.LongQueryTime = 1
	xdb.RegistryLogger(&dbLogger{
		name: "dbslowsql",
	})

	provider, err = opentelemetry.NewTracerProvider(global.AppName, global.Config)
	if err != nil {
		return err
	}

	return nil
}

func StopedHook(ctx sctx.Context) error {
	if provider != nil {
		return provider.Shutdown(ctx)
	}
	return nil
}
