package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Auth object that handle authentication
type Auth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	CookieDomain  string
	CookiePath    string
	CookieName    string
}

// user object (who is requesting authentication)
type jwtUser struct {
	ID        int    `json:"id"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
}

// token object that contains 2 tokens (refresh and access)
type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// claims object
type Claims struct {
	jwt.RegisteredClaims
}

// method of Auth object to generate TokenPair
func (j *Auth) GenerateTokenPair(user *jwtUser) (TokenPairs, error) {
	// Create a token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	// Assert token.Claims type as MapClaims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprintf("%s %s", user.Firstname, user.Lastname)
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = j.Audience
	claims["iss"] = j.Issuer
	claims["iat"] = time.Now().UTC().Unix()
	claims["typ"] = "JWT"

	// Set the expiry for JWT
	claims["exp"] = time.Now().UTC().Add(j.TokenExpiry).Unix()

	// Create a signed token
	signedAccessToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// Create a refresh token and set claims
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()

	// Set the expiry for the refresh token
	refreshTokenClaims["exp"] = time.Now().UTC().Add(j.RefreshExpiry).Unix()

	// Create signed refresh token
	signedRefreshToken, err := token.SignedString([]byte(j.Secret))

	if err != nil {
		return TokenPairs{}, err
	}

	// Create TokenPairs and populate with signed tokens
	var tokenPair = TokenPairs{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	// Return TokenPairs
	return tokenPair, nil
}
