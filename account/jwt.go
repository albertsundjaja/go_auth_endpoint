package account

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var EXPIRY = 24 * time.Hour
var JWT_SECRET = "super secret"

// Create JWT token for authentication and authorization
// Default expiry is 24 hours from now, to use default set the expiry to 0 (time.Time{})
func CreateJwt(username string, expiry time.Time) (string, error){
	var func_name="CreateJwt"

	zeroTime := time.Time{}
	if expiry == zeroTime {
		expiry = time.Now().Add(EXPIRY)
	}

	jwtClaims := jwt.MapClaims {
		"Username":  username,
		"ExpiresAt": expiry,
	}

	//create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	//sign token with our secret
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		fmt.Println(func_name, err)
		return "", err
	}

	return tokenString, nil
}

// validate and get jwt claims
func ValidateSigningAndGetJwtClaims(jwtToken string) (map[string]interface{}, error) {
	var func_name="ValidateSigningAndGetJwtClaims"

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(JWT_SECRET), nil})

	//token invalid signing
	if err != nil {
		fmt.Println(func_name, err)
		return map[string]interface{}{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		fmt.Println(func_name, "Unable to parse token claims")
		return map[string]interface{}{} ,fmt.Errorf("unable to parse token for claims")
	}

	return claims, nil
}

