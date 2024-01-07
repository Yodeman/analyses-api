package token

import (
	"testing"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/stretchr/testify/require"

	"github.com/yodeman/analyses-api/util"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomUser()
	duration := time.Minute

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomUser()
	duration := -time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
}

func TestInvalidPasetoToken(t *testing.T) {
	maker1, err := NewPasetoMaker(util.RandomString(32))
	payload, err := NewPasetoPayload(util.RandomUser(), time.Minute)
	require.NoError(t, err)

	pasetoToken, err := paseto.NewTokenFromClaimsJSON(payload, nil)
	require.NoError(t, err)
	token := pasetoToken.V4Encrypt(maker1.symmetricKey, nil)

	maker2, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	pasetoPayload, err := maker2.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, pasetoPayload)
}
