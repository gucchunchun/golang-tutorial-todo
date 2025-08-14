package utils

import (
	"fmt"

	"github.com/google/uuid"
)

func ParseUserUUID(input string) (uuid.UUID, error) {
	id, err := uuid.Parse(input)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID format: %w", err)
	}
	return id, nil
}
