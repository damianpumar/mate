package framework

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"time"
)

type SecureCookie struct {
	key []byte
}

func NewSecureCookie(key string) *SecureCookie {
	return &SecureCookie{
		key: adjustKeyLength([]byte(key)),
	}
}

func adjustKeyLength(key []byte) []byte {
	if len(key) < 16 {

		for len(key) < 16 {
			key = append(key, '0')
		}
		return key[:16]
	} else if len(key) < 24 {

		for len(key) < 24 {
			key = append(key, '0')
		}
		return key[:24]
	} else if len(key) < 32 {

		for len(key) < 32 {
			key = append(key, '0')
		}
		return key[:32]
	} else {

		return key[:32]
	}
}

func (sc *SecureCookie) SetEncryptedCookie(w http.ResponseWriter, name, value string, expiry time.Duration) error {
	encryptedValue, err := sc.encrypt(value)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    encryptedValue,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(expiry),
		Secure:   true,
	}

	http.SetCookie(w, cookie)
	return nil
}

func (sc *SecureCookie) GetEncryptedCookie(req *http.Request, name string) (string, error) {
	cookie, err := req.Cookie(name)

	if err != nil {
		return "", err
	}

	return sc.decrypt(cookie.Value)
}

func (sc *SecureCookie) HasCookie(req *http.Request, name string) bool {
	if _, err := sc.GetEncryptedCookie(req, name); err != nil {
		return false
	}

	return true
}

func (sc *SecureCookie) ClearCookie(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0), // Fecha en el pasado
		MaxAge:   -1,              // Invalida inmediatamente la cookie
		Secure:   true,            // Opcionalmente, puede establecerse solo para HTTPS
	}

	http.SetCookie(w, cookie)
}

func (sc *SecureCookie) encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(sc.key)
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

	encrypted := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.URLEncoding.EncodeToString(encrypted), nil
}

func (sc *SecureCookie) decrypt(cipherText string) (string, error) {
	decodedData, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(sc.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(decodedData) < nonceSize {
		return "", errors.New("invalid ciphertext")
	}

	nonce, cipherText := decodedData[:nonceSize], string(decodedData[nonceSize:])
	plainText, err := gcm.Open(nil, nonce, []byte(cipherText), nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
