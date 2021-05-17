package main

import (
	"time"
)

type LoginPostDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type MetricsPostDTO struct {
	AccountId string    `json:"account_id"`
	UserId    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

type MetricsGetResDTO struct {
	IotUserCount int `json:"iot_user_count"`
}
