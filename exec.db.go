package xdbmongodb

import (
	"context"

	contribxdb "github.com/zhiyunliu/glue/contrib/xdb"
	"github.com/zhiyunliu/glue/xdb"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongodb struct {
	cfg      *contribxdb.Setting
	client   *mongo.Client
	database *mongo.Database
}

func (db *mongodb) Query(ctx context.Context, sql string, input any) (data xdb.Rows, err error) {
	err = NotImplemented
	return
}

func (db *mongodb) Multi(ctx context.Context, sql string, input any) (data []xdb.Rows, err error) {
	err = NotImplemented
	return
}

func (db *mongodb) First(ctx context.Context, sql string, input any) (data xdb.Row, err error) {
	err = NotImplemented
	return
}

func (db *mongodb) Scalar(ctx context.Context, sql string, input any) (data interface{}, err error) {
	err = NotImplemented
	return
}

func (db *mongodb) Exec(ctx context.Context, sql string, input any) (r xdb.Result, err error) {
	err = NotImplemented
	return
}

func (db *mongodb) QueryAs(ctx context.Context, sql string, input any, result any) (err error) {
	return NotImplemented
}

func (db *mongodb) FirstAs(ctx context.Context, sql string, input any, result any) (err error) {
	return NotImplemented
}

func (db *mongodb) Begin() (trans xdb.ITrans, err error) {
	err = NotImplemented
	return
}

func (db *mongodb) Close() (err error) {
	if db.client != nil {
		return db.client.Disconnect(context.Background())
	}
	return
}

func (db *mongodb) GetImpl() (impl any) {
	return db.database
}

func (db *mongodb) Transaction(callback xdb.TransactionCallback) (err error) {
	return NotImplemented
}
