package db

import "time"

type Operation struct {
	Id         int64
	Name       string
	UserId     int64
	Scope      float64
	DateCreate time.Time
}
