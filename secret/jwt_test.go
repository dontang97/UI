package secret

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type _Suite struct {
	suite.Suite
}

func (s *_Suite) SetupSuite() {
	InitSecretKey(".")
}

func (s *_Suite) TearDownSuite() {
}

func (s *_Suite) SetupTest() {
}

func (s *_Suite) TearDownTest() {
}

func (s *_Suite) TestKey() {
	s.Equal(1702, len(priKey))
	s.Equal(460, len(pubKey))

	s.Equal(true, rsaPriKey != nil)
	s.Equal(true, rsaPubKey != nil)
}

func (s *_Suite) TestJWT() {
	acct := "kobe"
	token, err := CreateUserJWT(acct)
	s.Equal(nil, err)

	err = VerifyUserJWT(token, acct)
	s.Equal(nil, err)

	err = VerifyUserJWT(token, "james")
	s.Equal("The account is not matched", err.Error())
}

func TestRun(t *testing.T) {
	suite.Run(t, new(_Suite))
}
