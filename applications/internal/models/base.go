package models

import "time"

type Base struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Blockchain struct {
	EventId     int64  `json:"event_id"`
	BlockTime   int64  `json:"block_time"`
	BlockNumber int64  `json:"block_number"`
	LogIndex    int64  `json:"log_index"`
	TxHash      string `json:"tx_hash"`
}
