package secret

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	pubKeyFile = "/ui_rsa_pub.pem"
	priKeyFile = "/ui_rsa_pri.pem"

	validDuration time.Duration = time.Minute * 15
	//validDuration time.Duration = time.Second * 1
)

var (
	pubKey    []byte
	rsaPubKey *rsa.PublicKey

	priKey    []byte
	rsaPriKey *rsa.PrivateKey
)

type JWTErrorCode int
type JWTError struct {
	code JWTErrorCode
}

const (
	JWTUnknownError      JWTErrorCode = 0
	JWTNotActiveError    JWTErrorCode = 1
	JWTExpiredError      JWTErrorCode = 2
	JWTAcctNotMatchError JWTErrorCode = 3
	JWTNotAuthError      JWTErrorCode = 4
)

func (err *JWTError) Error() string {
	switch err.code {
	case JWTNotActiveError:
		return "JWT is not active yet"
	case JWTExpiredError:
		return "JWT is expired"
	case JWTAcctNotMatchError:
		return "The account is not matched"
	case JWTNotAuthError:
		return "Not authorized JWT"
	}

	return "Unknown JWT error"
}

const (
	JWTClaimFieldAcct = "acct"
	JWTClaimFieldAuth = "authorized"
	JWTClaimFieldExp  = "exp"
)

func InitSecretKey(keyDir string) {
	var err error
	if pubKey, err = ioutil.ReadFile(keyDir + pubKeyFile); err != nil {
		log.Fatal(err)
	}
	if rsaPubKey, err = jwt.ParseRSAPublicKeyFromPEM(pubKey); err != nil {
		log.Fatal(err)
	}
	if priKey, err = ioutil.ReadFile(keyDir + priKeyFile); err != nil {
		log.Fatal(err)
	}
	if rsaPriKey, err = jwt.ParseRSAPrivateKeyFromPEM(priKey); err != nil {
		log.Fatal(err)
	}
}

func CreateUserJWT(acct string) (string, error) {
	var err error

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims[JWTClaimFieldAuth] = true
	atClaims[JWTClaimFieldAcct] = acct
	atClaims[JWTClaimFieldExp] = time.Now().Add(validDuration).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims)
	if err != nil {
		return "", err
	}
	token, err := at.SignedString(rsaPriKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func VerifyUserJWT(tokenStr, acct string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return rsaPubKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return &JWTError{JWTNotActiveError}
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return &JWTError{JWTExpiredError}
			}

			return err
		}
	}

	if token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		s := claims[JWTClaimFieldAcct].(string)
		if s != acct {
			return &JWTError{JWTAcctNotMatchError}
		}
		if !claims[JWTClaimFieldAuth].(bool) {
			return &JWTError{JWTNotAuthError}
		}

		return nil
	}

	return &JWTError{JWTUnknownError}
}
