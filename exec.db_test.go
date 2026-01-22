package xdbmongodb_test

import (
	"testing"

	xdbmongodb "github.com/zhiyunliu/xdb-mongodb"
)

func TestTest(t *testing.T) {

	gotErr := xdbmongodb.Test()
	if gotErr != nil {
		t.Errorf("Test() failed: %v", gotErr)
	}
}
