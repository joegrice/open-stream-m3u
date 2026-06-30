package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

func EncodeToken(config any) (string, error) {
	data, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data), nil
}

func DecodeToken(token string) (map[string]any, error) {
	if strings.HasPrefix(token, "enc:") {
		return nil, errors.New("encrypted tokens require secret; use DecryptConfig")
	}

	// Try base64url first
	data, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		// Try standard base64
		data, err = base64.StdEncoding.DecodeString(token)
		if err != nil {
			// Try base64url with padding
			data, err = base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(token)
			if err != nil {
				return nil, fmt.Errorf("invalid base64: %w", err)
			}
		}
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return result, nil
}

func EncryptConfig(config any, secret string) (string, error) {
	if len(secret) < 16 {
		return "", errors.New("secret must be at least 16 characters")
	}

	data, err := json.Marshal(config)
	if err != nil {
		return "", err
	}

	key := deriveKey(secret)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return "enc:" + base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(ciphertext), nil
}

func DecryptConfig(token string, secret string) (map[string]any, error) {
	if !strings.HasPrefix(token, "enc:") {
		return nil, errors.New("token is not encrypted")
	}

	if len(secret) < 16 {
		return nil, errors.New("secret must be at least 16 characters")
	}

	encoded := strings.TrimPrefix(token, "enc:")
	ciphertext, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(encoded)
	if err != nil {
		ciphertext, err = base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return nil, err
		}
	}

	key := deriveKey(secret)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(plaintext, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func deriveKey(secret string) []byte {
	hash := sha256.Sum256([]byte(secret))
	return hash[:]
}
