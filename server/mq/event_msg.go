package mq

import (
	"encoding/json"
	"log"
)

type UserActivatedEvent struct {
	UserID       int   `json:"user_id"`
	ActivateTime int64 `json:"activate_time"`
}

func (u *UserActivatedEvent) ToBytes() []byte {
	ret, err := json.Marshal(u)
	if err != nil {
		log.Fatalf("Failed to marshal UserActivatedEvent: %v", err)
		return nil
	}
	return ret
}
