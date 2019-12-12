package models

import (
	"context"
	"log"
	"time"

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

type PageParams struct {
	Page     int ` json:"pageNo" form:"pageNo"`
	PageSize int ` json:"pageSize" form:"pageSize"`
}

type TimeRange struct {
	Start time.Time `json:"start,omitempty" form:"start,omitempty" time_format:"2006-01-02 15:04:05" time_utc:"1"`
	End   time.Time `json:"end,omitempty" form:"end,omitempty" time_format:"2006-01-02 15:04:05" time_utc:"1"`
}
type QueryParam struct {
	PageParams
	TimeRange
	Title string `json:"title" form:"title"`
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

func GetLogs(user string, param *QueryParam) []VisitLog {
	filter := bson.D{}
	param = parseQuery(param)
	if user != "" {
		filter = append(filter, bson.E{"user", user})
	}
	if param.Title != "" {
		filter = append(filter, bson.E{"title", primitive.Regex{Pattern: param.Title}})
	}
	filter = append(filter, bson.E{
		"begin_time", bson.M{
			"$gte": param.Start.Unix(),
			"$lte": param.End.Unix(),
		},
	})
	skip := int64((param.Page - 1) * param.PageSize)
	limit := int64(param.PageSize)
	option := &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
		Sort: bson.M{
			"begin_time": -1,
		},
	}
	cur, err := Db.Collection("visit").Find(context.TODO(), filter, option)
	if err != nil {
		log.Println(err)
		return []VisitLog{}
	}
	result := make([]VisitLog, 0)
	//m := make([]map[string]interface{}, 0)
	err = cur.All(context.TODO(), &result)
	if err != nil {
		log.Println(err.Error())
		return []VisitLog{}
	}
	return result
	//return result
}
func GetUserVisitCount(user string, param *QueryParam) int64 {
	filter := bson.D{
		{"user", user},
	}
	param = parseQuery(param)
	if param.Title != "" {
		filter = append(filter, bson.E{"title", param.Title})
	}
	filter = append(filter, bson.E{
		"begin_time", bson.M{
			"$gte": param.Start.Unix(),
			"$lte": param.End.Unix(),
		},
	})
	total, _ := Db.Collection("visit").CountDocuments(context.Background(), filter)
	return total
}

func parseQuery(q *QueryParam) *QueryParam {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
	if q.Start.IsZero() {
		q.Start = time.Now().AddDate(0, 0, -1)
	}
	if q.End.IsZero() {
		q.End = time.Now()
	}
	return q
}
