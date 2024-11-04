package services

import (
	"github.com/zhiyunliu/glue"
	"github.com/zhiyunliu/glue/context"
	"github.com/zhiyunliu/glue/server/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func BindAPI(srv *api.Server) {

	srv.Handle("/test", func(ctx context.Context) (res any) {

		val := ctx.Request().Query().Values()

		dbObj := glue.DB("mongodb")
		db := dbObj.GetImpl().(*mongo.Database)
		collection := db.Collection("x", options.Collection())

		if val.Get("c") == "1" {
			collection.InsertOne(ctx.Context(), val)
		}

		dbResult := collection.FindOne(ctx.Context(), val)
		if dbResult.Err() != nil {
			return dbResult.Err()
		}

		result := &DataRow{}
		err := dbResult.Decode(&result)
		if err != nil {
			return err
		}

		return result
	})

	srv.Handle("/cmd", func(ctx context.Context) (res any) {

		dbObj := glue.DB("mongodb")
		client := dbObj.GetImpl().(*mongo.Client)
		db := client.Database("test")

		filter := bson.D{{Key: "a", Value: 100}}

		dbRst := db.RunCommand(ctx.Context(), bson.D{
			{Key: "find", Value: "x"},
			{Key: "filter", Value: filter},
		})

		result := &DataRow{}
		err := dbRst.Decode(&result)
		if err != nil {
			return err
		}

		return result
	})
}

type DataRow struct {
	A string `json:"a"`
	B string
}
