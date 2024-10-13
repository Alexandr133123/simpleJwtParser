package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"
)

const secretKey = "123"

func CreateTokenHandler(response http.ResponseWriter, request *http.Request) {
	b := strings.Builder{}
	h := hmac.New(sha256.New, []byte(secretKey))

	header := "{\r\n  \"alg\": \"HS256\",\r\n  \"typ\": \"JWT\"\r\n}"
	payload := "{}"

	encodedHeader := base64.URLEncoding.EncodeToString([]byte(header))
	encodedPayload := base64.URLEncoding.EncodeToString([]byte(payload))

	b.WriteString(encodedHeader)
	b.WriteRune('.')
	b.WriteString(encodedPayload)

	h.Write([]byte(b.String()))
	dataHmac := h.Sum(nil)
	encodedHmac := base64.URLEncoding.EncodeToString(dataHmac)

	b.WriteRune('.')
	b.WriteString(encodedHmac)

	response.Write([]byte(b.String()))
}

func ValidateTokenHandler(response http.ResponseWriter, request *http.Request) {
	jwt := request.Header.Get("Authorization")

	segments := strings.Split(jwt, ".")
	if len(segments) != 3 {
		response.WriteHeader(401)
		return
	}

	signatureEncoded := segments[2]

	b := strings.Builder{}
	b.WriteString(segments[0])
	b.WriteRune('.')
	b.WriteString(segments[1])

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(b.String()))
	dataHmac := h.Sum(nil)
	encodedHmac := base64.URLEncoding.EncodeToString(dataHmac)

	if encodedHmac != signatureEncoded {
		response.WriteHeader(401)
		return
	}

	response.WriteHeader(200)
	response.Write([]byte("Token is valid\n"))
}

func helloHandler(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("Hello, world"))
}

func main() {
	m := http.NewServeMux()
	m.HandleFunc("POST /token", CreateTokenHandler)
	m.HandleFunc("GET /token", ValidateTokenHandler)
	m.HandleFunc("GET /hello", helloHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: m,
	}
	println("Listening port 8080")

	server.ListenAndServe()
}
