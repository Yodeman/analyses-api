package token

import (
    "encoding/json"
    "time"

    "aidanwoods.dev/go-paseto"
)

type PasetoMaker struct {
    symmetricKey paseto.V4SymmetricKey
}

func NewPasetoMaker(symmetricKey string) (*PasetoMaker, error) {
    v4SymmetricKey, err := paseto.V4SymmetricKeyFromBytes([]byte(symmetricKey))
    if err != nil {
        return nil, err
    }

    return &PasetoMaker{v4SymmetricKey}, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
    payload, err := NewPasetoPayload(username, duration)
    if err != nil {
        return "", err
    }

    token, err := paseto.NewTokenFromClaimsJSON(payload, nil)
    if err != nil {
        return "", err
    }

    return token.V4Encrypt(maker.symmetricKey, nil), nil

}

// VerifyToken checks if the token is valid or not
func (maker *PasetoMaker)VerifyToken(token string) (*PasetoPayload, error) {
    parser := paseto.NewParser()

    parsedToken, err := parser.ParseV4Local(maker.symmetricKey, token, nil)
    if err != nil {
        return nil, err
    }

    var payload PasetoPayload

    err = json.Unmarshal(parsedToken.ClaimsJSON(), &payload)
    if err != nil {
        return nil, err
    }

    return &payload, nil
}
