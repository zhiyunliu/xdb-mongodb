package main

import (
	"context"
	"fmt"
)

type dbLogger struct {
	name string
}

func (l dbLogger) Name() string {
	return l.name
}
func (l dbLogger) Log(ctx context.Context, elapsed int64, sql string, args ...interface{}) {
	fmt.Printf("%s:elapsed:%d, %s\n,args:%+v\n", l.name, elapsed, sql, args)
}
