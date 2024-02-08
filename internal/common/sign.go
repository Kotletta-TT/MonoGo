// Package common implements some utils
package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// CreateHMACSignature calculates the HMAC signature of a message using the provided key.
//
// Parameters:
// - key: the secret key used to generate the signature.
// - message: the message to calculate the signature for.
//
// Returns:
// - string: the hexadecimal representation of the HMAC signature.
// - error: an error if the calculation fails.
func CreateHMACSignature(key string, message []byte) (string, error) {
	h := hmac.New(sha256.New, []byte(key))

	_, err := h.Write(message)
	if err != nil {
		return "", err
	}

	signature := h.Sum(nil)
	return hex.EncodeToString(signature), nil
}

// VerifyHMACSignature verifies the HMAC signature of the given data using the provided key and signature.
//
// Parameters:
// - sign: The expected HMAC signature.
// - key: The key used for HMAC calculation.
// - data: The data to be verified.
//
// Returns:
// - error: An error if the calculated HMAC signature does not match the expected signature, or if there is an error during the verification process.
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
