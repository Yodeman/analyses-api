package token

import (
    "encoding/json"
    "time"

    "github.com/google/uuid"
)

// PasetoPayload contains the payload data of the PASETO token
type PasetoPayload struct {
    ID          string      `json:"jti,omitempty"`
    Username    string      `json:"preferred_username,omitempty"`
    IssuedAt    time.Time   `json:"iat,omitempty"`
    ExpiresAt   time.Time   `json:"exp,omitempty"`
}

// NewPasetoPayload creates a new PASETO payload with specific
// username and duration
func NewPasetoPayload(username string, duration time.Duration) ([]byte, error) {
    tokenID, err := uuid.NewRandom()
    if err != nil {
        return nil, err
    }

    p := &PasetoPayload{
        Username:   username,
        ID:         tokenID.String(),
        IssuedAt:   time.Now(),
        ExpiresAt:  time.Now().Add(duration),
    }

    payload, err := json.Marshal(&p)
    if err != nil {
        return nil, err
    }

    return payload, nil
}
