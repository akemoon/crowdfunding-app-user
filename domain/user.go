package domain

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	minUsernameLen = 3
	maxUsernameLen = 15
)

type CreateUserReq struct {
	UserID   uuid.UUID `json:"userID"`
	Username string    `json:"username"`
}

type CreateUserResp struct {
	UserID uuid.UUID `json:"userID"`
}

type UpdateUsername struct {
}

type UpdateDescription struct {
}

type User struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Description string    `json:"description"`
}

func ValidateUsernameLen(username string) error {
	if len(username) < minUsernameLen || len(username) > maxUsernameLen {
		return fmt.Errorf("username len must be between %d and %d", minUsernameLen, maxUsernameLen)
	}
	return nil
}
