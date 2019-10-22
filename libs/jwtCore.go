/*
 * Copyright (c) 2019. dvnlabs.ml and Buddylang
 * Davin Alfarizky Putra Basudewa <dbasudewa@gmail.com>
 * This project is for Buddylang Api for core functional.Dont expose to public
 * Please kept this source private.
 */
package libs

import (
	conf "absensi-server/config"
	"absensi-server/core/master/user/model"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	jwt.StandardClaims
	Data interface{}
}

func NewToken(data interface{}) (bool, string) {
	//expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: 0,
		},
		Data: data,
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(conf.TokenSecretEncoded())
	if err != nil {
		println(err.Error())
		return true, ""
	}
	return false, tokenString
}
func VerifyToken(token string) (bool, interface{}, string) {
	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return conf.TokenSecretEncoded(), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, model.LoginData{}, err.Error()
		}
		return false, model.LoginData{}, err.Error()
	}
	if !tkn.Valid {
		return false, model.LoginData{}, "INVALID TOKEN"
	}
	return true, claims.Data, "SUCCESS Parsing"
}
