package auth

import (
	"banner-service/internal/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/jellydator/ttlcache/v3"
	"time"
)

type TokenProvider struct {
	Cache      *ttlcache.Cache[string, models.UserResources]
	PrivateKey string
	PublicKey  string
	TokenTTL   time.Duration
}

func NewTokenProvider(cache *ttlcache.Cache[string, models.UserResources], privateKey, publicKey string, tokenTTL time.Duration) TokenProvider {
	return TokenProvider{
		Cache:      cache,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		TokenTTL:   tokenTTL,
	}
}

func (tp *TokenProvider) CreateToken(ttl time.Duration, username string, payload models.UserResources, privateKey string) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["sub"] = username
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(ttl).Unix()
	claims["resources"] = payload

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)

	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func (tp *TokenProvider) ValidateToken(token string, publicKey string) (models.UserResources, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return models.UserResources{}, fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return models.UserResources{}, fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return models.UserResources{}, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return models.UserResources{}, fmt.Errorf("validate: invalid token")
	}

	jsonResource, err := json.Marshal(claims["resources"])
	if err != nil {
		return models.UserResources{}, fmt.Errorf("validate: %w", err)
	}

	var resources models.UserResources

	if err = json.Unmarshal(jsonResource, &resources); err != nil {
		return models.UserResources{}, fmt.Errorf("validate: %w", err)
	}

	return resources, nil
}
