package xdbmongodb

import (
	"context"
	"sync"

	"github.com/zhiyunliu/glue/global"
	"github.com/zhiyunliu/glue/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	onceLock sync.Once
	meter    *Metrics
)

type Metrics struct {
	RequestCounter metrics.Int64Counter `metric:"req_total"  `
}

func InitMetrics() (err error) {
	onceLock.Do(func() {
		meter = &Metrics{}
		factory := metrics.NewFactory(otel.GetMeterProvider(), ScopeName)
		err = metrics.Init(meter, factory)
	})
	return err
}

func (m *Metrics) IncRequest(connName string) {
	if m.RequestCounter == nil {
		return
	}
	m.RequestCounter.Add(context.Background(), 1, metric.WithAttributes(attribute.String("conn", connName), attribute.String("srv", global.AppName)))
}
