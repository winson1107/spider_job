package models

import (
	"context"
)

type VisitLog struct {
	BeginTime int64  `json:"begin_time" bson:"begin_time"`
	EndTime   int64  `json:"end_time" bson:"end_time" `
	Host      string `json:"host" bson:"host"`
	Url       string `json:"url" bson:"url"`
	User      string `json:"uid" bson:"user"`
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
