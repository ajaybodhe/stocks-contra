package coreStructures

import (
	"time"
)

type NseBhavRecord struct {
	RecordDate     time.Time
	PrevClosePrice float32
	OpenPrice      float32
	HighPrice      float32
	LowPrice       float32
	LastPrice      float32
	ClosePrice     float32
	AvgPrice       float32
	TtlTrdQnty     int
	DelivQty       int
	DelivPer       float32
}
type NseBhavData struct {
	Symbol     string
	BhavRecord []NseBhavRecord
}

func NewNseBhavData(daysCount int) *NseBhavData {
	var c NseBhavData
	c.BhavRecord = make([]NseBhavRecord, daysCount)
	return &c
}
