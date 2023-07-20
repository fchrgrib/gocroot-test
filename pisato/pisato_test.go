package pisato

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker("qZertCuiOpLMNFDsChJkIoLAnHpkCDEf")

	require.NoError(t, err)

	email := "Fahrian.Afdholi@gmail.com"
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(email, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)

	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.Id)
	require.Equal(t, email, payload.UserName)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker("qZertCuiOpLMNFDsChJkIoLAnHpkCDEf")
	require.NoError(t, err)

	token, err := maker.CreateToken("Fahrian.Afdholi@gmail.com", -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, errors.New("token is expired").Error())
	require.Nil(t, payload)
}

func TestLengthOfKey(t *testing.T) {
	maker, err := NewPasetoMaker("qZertCuiOpLMNFDsChJkIoLAnHpkCDE")

	require.Error(t, err)
	require.EqualError(t, err, errors.New("error because symmetric doesn't contain 32 length").Error())
	require.Nil(t, maker)
}
