package api

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

func NewHTTPSignature(key []byte) (map[string]interface{}, error) {
	headers := make(map[string]interface{})

	contentType := "application/json"
	encodedAuth := base64.StdEncoding.EncodeToString(key)

	headers["Accept"] = contentType
	headers["Content-Type"] = contentType
	headers["Authorization"] = fmt.Sprintf("Basic %s", encodedAuth)

	return headers, nil
}
