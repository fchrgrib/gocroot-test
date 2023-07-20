package pisato

import (
	"errors"
	"github.com/aead/chacha20poly1305"
	"github.com/gocroot/gocroot/handler"
	"github.com/o1egl/paseto"
	"time"
)

type PasetoMaker struct {
	paseto *paseto.V2
	key    []byte
}

func NewPasetoMaker(key string) (handler.Maker, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, errors.New("error because symmetric doesn't contain 32 length")
	}

	maker := &PasetoMaker{
		paseto: paseto.NewV2(),
		key:    []byte(key),
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := handler.NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	encrypt, err := maker.paseto.Encrypt(maker.key, payload, nil)
	if err != nil {
		return "", err
	}

	return encrypt, nil
}

func (maker *PasetoMaker) VerifyToken(token string) (*handler.Payload, error) {
	payload := &handler.Payload{}

	if err := maker.paseto.Decrypt(token, maker.key, payload, nil); err != nil {
		return nil, err
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}

	return payload, nil
}
