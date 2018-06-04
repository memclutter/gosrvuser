package main

import "time"

// Common response format
type Response struct {
	Success bool        `json:"success"`
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponseDataStatus struct {
	Time   time.Time `json:"time"`
	Health struct {
		Db   bool `json:"db"`
		Amqp bool `json:"amqp"`
	} `json:"health"`
}
