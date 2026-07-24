package xdbmongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/event"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	cmap "github.com/orcaman/concurrent-map/v2"
	"go.mongodb.org/mongo-driver/bson"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type spanKey struct {
	ConnectionID string
	RequestID    int64
}

type cacheItem struct {
	spanName string
	span     trace.Span
	cmd      bson.Raw
}

type monitor struct {
	TracerProvider           trace.TracerProvider
	Tracer                   trace.Tracer
	CommandAttributeDisabled bool

	cfg      *monitorConfig
	cmdCache cmap.ConcurrentMap[spanKey, cacheItem]
}

func NewMonitor(cfg *monitorConfig) *event.CommandMonitor {

	shardCnt := int64(cmap.SHARD_COUNT)

	m := &monitor{
		TracerProvider: otel.GetTracerProvider(),
		cfg:            cfg,
		cmdCache: cmap.NewWithCustomShardingFunction[spanKey, cacheItem](func(key spanKey) uint32 {
			return uint32(key.RequestID % shardCnt)
		}),
	}
	m.Tracer = m.TracerProvider.Tracer("XDB.MONGODB")
	return &event.CommandMonitor{
		Started:   m.Started,
		Succeeded: m.Succeeded,
		Failed:    m.Failed,
	}

}

func (m *monitor) Started(ctx context.Context, evt *event.CommandStartedEvent) {
	var spanName string
	collection, err := extractCollection(evt)
	if err == nil && collection != "" {
		spanName = collection + "."
	}
	spanName += evt.CommandName

	meter.Incr(m.cfg.ConnName, spanName)

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(m.commandStartedTraceAttrs(evt, collection)...),
	}
	_, span := m.Tracer.Start(ctx, spanName, opts...)
	m.cmdCache.Set(
		spanKey{
			RequestID:    evt.RequestID,
			ConnectionID: evt.ConnectionID,
		},
		cacheItem{
			spanName: spanName,
			span:     span,
			cmd:      evt.Command,
		},
	)
}

func (m *monitor) Succeeded(ctx context.Context, evt *event.CommandSucceededEvent) {
	m.Finished(ctx, &evt.CommandFinishedEvent, nil)
}

func (m *monitor) Failed(ctx context.Context, evt *event.CommandFailedEvent) {
	m.Finished(ctx, &evt.CommandFinishedEvent, fmt.Errorf("%s", evt.Failure))
}

func (m *monitor) Finished(ctx context.Context, evt *event.CommandFinishedEvent, err error) {
	key := spanKey{
		ConnectionID: evt.ConnectionID,
		RequestID:    evt.RequestID,
	}
	item, ok := m.cmdCache.Get(key)
	if ok {
		m.cmdCache.Remove(key)
	} else {
		return
	}

	if err != nil {
		item.span.SetStatus(codes.Error, err.Error())
	}
	item.span.End()
	meter.Decr(m.cfg.ConnName, item.spanName)

	m.printSlowQuery(ctx, evt.RequestID, evt.Duration, evt.CommandName, item.cmd)

}

func extractCollection(evt *event.CommandStartedEvent) (connName string, err error) {
	elt, err := evt.Command.IndexErr(0)
	if err != nil {
		return
	}
	var key string
	if key, err = elt.KeyErr(); err == nil && key == evt.CommandName {
		var v bson.RawValue
		if v, err = elt.ValueErr(); err != nil || v.Type != bson.TypeString {
			return
		}
		return v.StringValue(), nil
	}
	err = errors.New("collection name not found")
	return
}

func (m *monitor) commandStartedTraceAttrs(evt *event.CommandStartedEvent, connName string) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.DBSystemMongoDB,
		semconv.DBOperationName(evt.CommandName),
		semconv.DBNamespace(evt.DatabaseName),
		semconv.NetworkTransportTCP,
		attribute.String("db.conn.name", m.cfg.ConnName),
		attribute.String("db.command", evt.Command.String()),
	}
	if connName != "" {
		attrs = append(attrs, semconv.DBCollectionName(connName))
	}

	return attrs
}

func (m *monitor) printSlowQuery(ctx context.Context, requestId int64, timeRange time.Duration, query string, cmd bson.Raw) {
	if !m.cfg.ShowQueryLog {
		return
	}
	if m.cfg.logger == nil {
		return
	}

	if m.cfg.slowThreshold <= 0 || timeRange < m.cfg.slowThreshold {
		return
	}
	m.cfg.logger.Log(ctx, timeRange.Milliseconds(), m.cfg.ConnName, fmt.Sprintf("[%s][%d]%s", m.cfg.ConnName, requestId, query), cmd.String())
}
