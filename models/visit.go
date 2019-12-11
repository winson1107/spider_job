package models

import (
	"context"

	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VisitLog struct {
	BeginTime int64              `json:"begin_time" bson:"begin_time" form:"begin_time" binding:"required"`
	EndTime   int64              `json:"end_time" bson:"end_time" form:"end_time" binding:"required" `
	Host      string             `json:"host" bson:"host"`
	Url       string             `json:"url" bson:"url" form:"url" binding:"required"`
	User      string             `json:"uid" bson:"user"`
	Title     string             `json:"title" bson:"title" form:"title"`
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
}

func InsertLogs(visitors []VisitLog) error {
	collect := Db.Collection("visit")
	rows := []interface{}{}
	for _, v := range visitors {
		rows = append(rows, v)
	}
	_, err := collect.InsertMany(context.Background(), rows)
	if err != nil {
		return err
	}
	return nil
}

func GetLogs(user string, page, limit int64) []VisitLog {
	filter := bson.D{
		{"user", user},
	}
	skip := int64((page - 1) * limit)
	option := &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
		Sort: bson.M{
			"begin_time": -1,
		},
	}
	cur, _ := Db.Collection("visit").Find(context.TODO(), filter, option)
	result := make([]VisitLog, 0)
	//m := make([]map[string]interface{}, 0)
	err := cur.All(context.TODO(), &result)
	if err != nil {
		log.Println(err.Error())
		return []VisitLog{}
	}
	return result
	//return result
}
func GetUserVisitCount(user string) int64 {
	filter := bson.D{
		{"user", user},
	}
	total, _ := Db.Collection("visit").CountDocuments(context.Background(), filter)
	return total
}
