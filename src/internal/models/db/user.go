package db

type User struct {
	Id         int64
	Balance    float64
	Operations []Operation
}
