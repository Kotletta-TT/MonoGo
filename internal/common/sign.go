package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

func CreateHMACSignature(key string, message []byte) (string, error) {
	h := hmac.New(sha256.New, []byte(key))

	_, err := h.Write(message)
	if err != nil {
		return "", err
	}

	signature := h.Sum(nil)
	return hex.EncodeToString(signature), nil
}

func VerifyHMACSignature(sign, key string, data []byte) error {
	h := hmac.New(sha256.New, []byte(key))

	_, err := h.Write(data)
	if err != nil {
		return err
	}

	calcSign := hex.EncodeToString(h.Sum(nil))

	if calcSign != sign {
		return errors.New("invalid signature")
	}

	return nil
}
