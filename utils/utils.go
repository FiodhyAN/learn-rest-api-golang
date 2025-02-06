package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FiodhyAN/learn-rest-api-golang/config"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func init() {
	Validate.RegisterValidation("password", validatePassword)
}

type JSONResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}

func ParseJSON(r *http.Request, value any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}

	return json.NewDecoder(r.Body).Decode(value)
}

func WriteJSON(w http.ResponseWriter, status int, message string, value interface{}, errors interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	response := JSONResponse{
		Status:  status,
		Message: message,
		Data:    value,
		Errors:  errors,
	}
	return json.NewEncoder(w).Encode(response)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, "error", nil, err.Error())
}

var encryptKey = []byte(config.Envs.EncryptKey)
var encryptIv = []byte(config.Envs.EncryptIv)

func EncryptText(plainText string) (string, error) {
	block, err := aes.NewCipher(encryptKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	plainTextPadded := PKCS7Padding([]byte(plainText), aes.BlockSize)
	ciphertext := make([]byte, len(plainTextPadded))

	mode := cipher.NewCBCEncrypter(block, encryptIv)
	mode.CryptBlocks(ciphertext, plainTextPadded)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptText decrypts the base64 encoded ciphertext
func DecryptText(encryptedText string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	block, err := aes.NewCipher(encryptKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("invalid ciphertext length")
	}

	mode := cipher.NewCBCDecrypter(block, encryptIv)
	mode.CryptBlocks(ciphertext, ciphertext)

	plainText, err := PKCS7UnPadding(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to unpad data: %w", err)
	}

	return string(plainText), nil
}

// PKCS7Padding applies padding
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7UnPadding removes padding
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("invalid padding: empty data")
	}

	unpadding := int(origData[length-1])
	if unpadding <= 0 || unpadding > length {
		log.Println("Invalid padding size:", unpadding, "for data length:", length)
		return nil, fmt.Errorf("invalid padding size: %d", unpadding)
	}

	return origData[:(length - unpadding)], nil
}

func FormatDate(date time.Time) (string, error) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return "", err
	}
	expirationDateInLocation := date.In(location)
	formattedDate := expirationDateInLocation.Format("Mon Jan 2 2006 15:04:05")
	timeZoneName := "Western Indonesia Time"
	finalFormattedDate := fmt.Sprintf("%s (%s)", formattedDate, timeZoneName)

	return finalFormattedDate, nil
}
