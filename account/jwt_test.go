package account

import (
	"testing"
	"time"
)

var dummyUsername = "username"
var expectedJwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOiIyMjAwLTAxLTAxVDE1OjAwOjAwWiIsIlVzZXJuYW1lIjoidXNlcm5hbWUifQ.v0ACp4SXlnWx6Zq23WQGqxTYtucIB7EuWuie9PoOMVU"

func TestCreateJwt(t *testing.T) {
	timeExpiry, _ := time.Parse("2006-01-02 15:04:05", "2200-01-01 15:00:00")
	jwt, err := CreateJwt(dummyUsername, timeExpiry)

	if err != nil {
		t.Errorf("Error creating jwt, expected error not to have occured")
	}

	if jwt != expectedJwt {
		t.Errorf("Jwt created is not the same, expected %s, got %s", expectedJwt, jwt)
	}
}

func TestValidateSigningAndGetJwtClaims(t *testing.T) {
	claims, err := ValidateSigningAndGetJwtClaims(expectedJwt)

	if err != nil {
		t.Errorf("Error checking jwt signing, expected error not to have occured, got %s", err)
	}

	expectedTimeExpiry, _ := time.Parse("2006-01-02 15:04:05", "2200-01-01 15:00:00")
	claimsExpiry, _ := time.Parse("2006-01-02T15:04:05Z", claims["ExpiresAt"].(string))
	if (claims["Username"] != dummyUsername || claimsExpiry != expectedTimeExpiry) {
		t.Errorf("Error in claims, expected username %s, expiry %s. Got %s, %s", dummyUsername, expectedTimeExpiry, claims["Username"], claims["ExpiresAt"])
	}
}