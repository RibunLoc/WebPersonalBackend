package util

import (
	"encoding/json"
	"fmt"
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

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	t, err := time.ParseInLocation(layout, str, loc)
	if err != nil {
		return err
	}
	*ct = CustomTime(t)
	return nil
}

// Lưu UTC vào MongoDB
func (ct CustomTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	t := time.Time(ct).UTC()
	millis := t.UnixMilli()
	buf := bsoncore.AppendDateTime(nil, millis)
	return bsontype.DateTime, buf, nil
}

// Khi đọc ra từ MongoDB, convert sang VN
func (ct *CustomTime) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	millis, _, ok := bsoncore.ReadDateTime(data)
	if !ok {
		return fmt.Errorf("failed to read bson datetime")
	}
	*ct = CustomTime(time.UnixMilli(millis).UTC())
	return nil
}

func (ct CustomTime) String() string {
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	vnTime := time.Time(ct).In(loc)
	return vnTime.Format(layout)
}
