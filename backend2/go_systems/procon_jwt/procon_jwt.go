package procon_jwt

import (
	"crypto/rsa"
	jwtgo "github.com/dgrijalva/jwt-go"
	"time"
)

func GenerateJWT(privatekeyfile *rsa.PrivateKey)(string,error) {
	token:= jwtgo.New(jwtgo.SigningMethodRS256)
	in10m := time.Now().Add(time.Duration(30)*time.Minute).Unix()
	token.Claims =jwtgo.MapClaims{
		"iss" : "procon.com",
		"aud" : "void.procon.com",
		"exp" : in10m,
		"iti" : "Unique",
		"iat" : time.Now().Unix(),
		"nbf" : 2,
		"sub" : "subject",
		"scopes" : "api:read",
	}
	tokenString,err := token.SignedString(privatekeyfile)
	if err!=nil {return "",err} else { return tokenString,nil }
}

func ValidateJWT(publickeyfile *rsa.PublicKey,jwt string)(bool, error){
	token,err := jwtgo.Parse(jwt,func(token *jwtgo.Token)(interface{}, error){
		return publickeyfile, nil
	})

	if err != nil { return false, err } else if (token.Valid && err==nil ) { return true,nil }


	return false,err;


}
