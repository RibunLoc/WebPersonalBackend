package util

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

type CustomTime time.Time

const layout = "2006-01-02 15:04:05"

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	vnTime := time.Time(ct).In(loc)
	return json.Marshal(vnTime.Format(layout))
}

func (ct CustomTime) String() string {
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	return time.Time(ct).In(loc).Format(layout)
}

func (ct CustomTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	t := time.Time(ct).UTC()
	return bsontype.DateTime, bsoncore.AppendDateTime(nil, t.UnixMilli()), nil
}
