package jwt

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"
)

type JWTManager interface {
	GenerateToken(payload interface{}) string
	AuthorizeToken(token string, p any) error
}

type SHA256JWTManager struct {
	ExpireDuration time.Duration
	secret         string
	providerName   string
}

// type tokenHeader struct {
// 	Typ       string
// 	Algorithm string
// }

// type TokenBody struct {
// 	Expiration time.Time
// 	Payload    interface{}
// 	Provider   string
// }

func (jm *SHA256JWTManager) GenerateToken(payload interface{}) string {
	header := map[string]string{
		"Typ":       "JWT",
		"Algorithm": "SHA256",
	}

	headerdata, err := json.Marshal(header)
	if err != nil {
		log.Panic(err)
	}

	body := map[string]interface{}{
		"Expiration": time.Now().Add(jm.ExpireDuration),
		"Payload":    payload,
		"Provider":   jm.providerName,
	}

	bodydata, err := json.Marshal(body)

	tokendata := base64.URLEncoding.EncodeToString(headerdata) + "." + base64.URLEncoding.EncodeToString(bodydata)

	if err != nil {
		log.Panic(err)
	}

	hasher := sha256.New()
	hasher.Write([]byte(tokendata + jm.secret))
	signature := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return tokendata + "." + signature

}

func (jm *SHA256JWTManager) AuthorizeToken(token string, p any) error {
	token_segments := strings.Split(token, ".")

	tokendata := token_segments[0] + "." + token_segments[1]

	hasher := sha256.New()
	hasher.Write([]byte(tokendata + jm.secret))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	if hash == token_segments[2] {
		// var tokenBody TokenBody
		var bodydata struct {
			Expiration time.Time
			Payload    json.RawMessage
			Provider   string
		}
		data, err := base64.URLEncoding.DecodeString(token_segments[1])
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(data, &bodydata)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(bodydata.Payload, p)
		if err != nil {
			log.Fatal(err)
		}

		if time.Now().After(bodydata.Expiration) {
			return errors.New("token Expired")
		}

		return nil
	}
	return errors.New("signature Cant be verified")
}

func NewJWTManager(expirationDuration time.Duration, secret string, providerName string) JWTManager {
	return &SHA256JWTManager{
		secret:         secret,
		providerName:   providerName,
		ExpireDuration: expirationDuration,
	}
}
