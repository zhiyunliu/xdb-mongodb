package xdbmongodb

import (
	"context"
	"sync"

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
	RequestCounter metrics.Int64UpDownCounter `metric:"db_cur_proc"  `
}

func InitMetrics() (err error) {
	onceLock.Do(func() {
		meter = &Metrics{}
		factory := metrics.NewFactory(otel.GetMeterProvider(), ScopeName)
		err = metrics.Init(meter, factory)
	})
	return err
}

func (m *Metrics) Incr(connName, spanName string) {
	if m.RequestCounter == nil {
		return
	}
	m.RequestCounter.Add(context.Background(), 1, metric.WithAttributes(attribute.String("dbtype", "mongo"), attribute.String("conn", connName), attribute.String("span", spanName)))
}

func (m *Metrics) Decr(connName, spanName string) {
	if m.RequestCounter == nil {
		return
	}
	m.RequestCounter.Add(context.Background(), -1, metric.WithAttributes(attribute.String("dbtype", "mongo"), attribute.String("conn", connName), attribute.String("span", spanName)))
}
