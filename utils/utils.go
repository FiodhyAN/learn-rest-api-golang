package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

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

func EncryptText(plainText string) (string, error) {
	var plainTextBlock []byte
	length := len(plainText)

	if length%16 != 0 {
		extendBlock := 16 - (length % 16)
		plainTextBlock = make([]byte, length+extendBlock)
		copy(plainTextBlock[length:], bytes.Repeat([]byte{uint8(extendBlock)}, extendBlock))
	} else {
		plainTextBlock = make([]byte, length)
	}

	copy(plainTextBlock, plainText)
	block, err := aes.NewCipher([]byte(config.Envs.EncryptKey))

	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, []byte(config.Envs.EncryptIv))
	mode.CryptBlocks(ciphertext, plainTextBlock)

	str := base64.StdEncoding.EncodeToString(ciphertext)

	return str, nil
}

func DecryptText(encryptText string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptText)

	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(config.Envs.EncryptKey))

	if err != nil {
		return nil, err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("block size cant be zero")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(config.Envs.EncryptIv))
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = PKCS5UnPadding(ciphertext)

	return ciphertext, nil
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	return src[:(length - unpadding)]
}
