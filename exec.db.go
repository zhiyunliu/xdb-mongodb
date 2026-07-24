package xdbmongodb

import (
	"context"
	"fmt"

	contribxdb "github.com/zhiyunliu/glue/contrib/xdb"
	"github.com/zhiyunliu/glue/xdb"
	"go.mongodb.org/mongo-driver/mongo"
)

// Test 函数用于验证当前 mongodb 结构体是否实现了 xdb.IDB 接口。
// 该测试通过创建临时 mongodb 实例并进行接口类型断言来检测接口实现情况。
//
// 返回值:
//
//	error:
//	  - 当成功实现 xdb.IDB 接口时返回 nil
//	  - 当未实现接口时返回包含错误信息的 error 对象
func Test() error {
	// 创建临时 mongodb 实例并赋值给空接口变量
	var tmpdb any = &mongodb{}

	// 尝试将临时对象断言为 xdb.IDB 接口
	// 成功断言说明已实现接口，返回 nil
	if _, ok := tmpdb.(xdb.IDB); ok {
		return nil
	}

	// 断言失败说明未实现接口，返回错误信息
	return fmt.Errorf("mongodb.Test():xdb.IDB not implemented")
}

var _ xdb.IDB = (*mongodb)(nil)

type mongodb struct {
	connName string
	cfg      *contribxdb.Setting
	client   *mongo.Client
	database *mongo.Database
}

func (db *mongodb) Query(ctx context.Context, sql string, input any, opts ...xdb.TemplateOption) (data xdb.Rows, err error) {
	err = NotImplemented
	return
}

func (db *mongodb) Multi(ctx context.Context, sql string, input any, opts ...xdb.TemplateOption) (data []xdb.Rows, err error) {
	err = NotImplemented
	return
}

func (db *mongodb) First(ctx context.Context, sql string, input any, opts ...xdb.TemplateOption) (data xdb.Row, err error) {
	err = NotImplemented
	return
}

func (db *mongodb) Scalar(ctx context.Context, sql string, input any, opts ...xdb.TemplateOption) (data interface{}, err error) {
	err = NotImplemented
	return
}

func (db *mongodb) Exec(ctx context.Context, sql string, input any, opts ...xdb.TemplateOption) (r xdb.Result, err error) {
	err = NotImplemented
	return
}

func (db *mongodb) QueryAs(ctx context.Context, sql string, input any, result any, opts ...xdb.TemplateOption) (err error) {
	return NotImplemented
}

func (db *mongodb) FirstAs(ctx context.Context, sql string, input any, result any, opts ...xdb.TemplateOption) (err error) {
	return NotImplemented
}

func (db *mongodb) Begin() (trans xdb.ITrans, err error) {
	err = NotImplemented
	return
}

func (db *mongodb) BeginTx(context.Context) (trans xdb.ITrans, err error) {
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
	db.IncRequest()
	return db.database
}

func (db *mongodb) Transaction(ctx context.Context, callback xdb.TransactionCallback) (err error) {
	return NotImplemented
}

func (db *mongodb) IncRequest() {
	meter.IncRequest(db.connName)
}
