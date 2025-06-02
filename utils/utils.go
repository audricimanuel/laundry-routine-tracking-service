package utils

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func TimeNow() time.Time {
	return time.Now()
}

func ConvertStrToInt(number string, defaultResult int) int {
	result, err := strconv.Atoi(number)
	if err != nil {
		return defaultResult
	}
	return result
}

func UnescapeString(text string) string {
	res, _ := url.QueryUnescape(strings.TrimSpace(strings.ToLower(text)))
	return res
}

func DecodeHttpResponse(data io.ReadCloser, dest interface{}) {
	json.NewDecoder(data).Decode(dest)
}

func ObjToByte(data interface{}) []byte {
	v, _ := json.Marshal(data)
	return v
}

func ConvertStructToMap(v any) map[string]interface{} {
	var result map[string]interface{}
	vByte, _ := json.Marshal(v)
	json.Unmarshal(vByte, &result)

	return result
}

func ConvertMapToStruct(data map[string]interface{}, target interface{}) error {
	byteData, _ := json.Marshal(data)
	if err := json.Unmarshal(byteData, target); err != nil {
		return err
	}
	return nil
}

func PtrToValue[T any](valuePtr *T) T {
	var v T
	if valuePtr != nil {
		v = *valuePtr
	}
	return v
}

func CustomError(message string) error {
	return errors.New(message)
}

func ConvertErrorToMap(key string, err error) map[string]string {
	return map[string]string{key: err.Error()}
}

func GenerateSlug(text string) string {
	result := strings.ToLower(strings.TrimSpace(text))
	result = strings.Replace(result, " ", "-", -1)
	return result
}

func GenerateUUID() uuid.UUID {
	return uuid.New()
}

func GenerateCleanUUID() string {
	traceId := GenerateUUID()
	cleanTraceId := strings.Replace(traceId.String(), "-", "", -1)
	return cleanTraceId
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a plain-text password with a hashed password.
// Returns true if they match, false otherwise.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
