package models

import (
	"context"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollectItem struct {
	Url     string `json:"url"`
	Created string `json:"created"`
	Title   string `json:"title"`
}

func GetCollects(page, pageSize int64) []CollectItem {
	filter := bson.D{}
	skip := int64((page - 1) * pageSize)
	option := &options.FindOptions{
		Limit: &pageSize,
		Skip:  &skip,
		Sort: bson.M{
			"created": -1,
		},
	}
	cur, _ := Db.Collection("collect").Find(context.TODO(), filter, option)
	result := make([]CollectItem, 0)
	_ = cur.All(context.TODO(), &result)
	return result
}

func AddCollectUrl(item *CollectItem) bool {
	item.Created = time.Now().Format("2006-01-02 15:04:05")
	res, err := Db.Collection("collect").InsertOne(context.TODO(), *item)
	if err != nil {
		return false
	}
	return res == nil
}

func QueryCollect(start, end time.Time) []CollectItem {
	filter := bson.D{
		{"created", bson.M{
			"$gt": start.Format("2006-01-02 15:04:05"),
			"$lt": end.Format("2006-01-02 15:04:05"),
		}},
	}
	option := &options.FindOptions{
		Sort: bson.M{
			"created": -1,
		},
	}
	cur, _ := Db.Collection("collect").Find(context.TODO(), filter, option)
	result := make([]CollectItem, 0)
	_ = cur.All(context.TODO(), &result)
	return result
}
