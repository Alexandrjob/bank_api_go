package dto

import (
	db2 "bank_api/src/internal/models/db"
)

type UpdateInfo struct {
	User      db2.User
	Operation db2.Operation
}
