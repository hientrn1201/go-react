package main

import (
	"fmt"
	"net/http"
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

// return a pointer to a http Cookie
func (j *Auth) GetRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Path:     j.CookiePath, // specifies the subset of URLs within the domain that the cookie applies to
		Value:    refreshToken,
		Expires:  time.Now().UTC().Add(j.RefreshExpiry), // provides an absolute duration (use for older browsers)
		MaxAge:   int(j.RefreshExpiry.Seconds()),        // provides a relative duration of expiration
		SameSite: http.SameSiteStrictMode,
		Domain:   j.CookieDomain, // specifies the domain that the cookie is associated with
		HttpOnly: true,           // deny access from JS
		Secure:   true,           // Cookie can only be sent over HTTPS connection
	}
}

// function that creates a cookie that is immediately expired - delete a cookie
func (j *Auth) GetExpiredRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Path:     j.CookiePath,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		Domain:   j.CookieDomain,
		HttpOnly: true,
		Secure:   true,
	}
}
