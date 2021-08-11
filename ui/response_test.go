package ui_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dontang97/ui/pg"
	"github.com/dontang97/ui/ui"
	"github.com/stretchr/testify/suite"
)

type _respSuite struct {
	suite.Suite
}

func (s *_respSuite) SetupSuite() {
}

func (s *_respSuite) TearDownSuite() {
}

func (s *_respSuite) SetupTest() {
}

func (s *_respSuite) TearDownTest() {
}

func (s *_respSuite) TestJsonIdent() {
	resp := ui.Response{
		Info: ui.Info{
			Status:  ui.StatusOK,
			Message: "This is testing message.",
		},
		Data: map[string]interface{}{
			"fullname": "James",
			"acct":     []string{"kobe", "jason"}},
	}

	_, err := resp.JsonIdent()
	//fmt.Println(string(b))
	s.Equal(nil, err)
}

func (s *_respSuite) TestWriteJsonResponse() {

	status := ui.Status(-1)
	data := pg.User{
		Acct:       "acct1",
		Pwd:        "pwd1",
		Fullname:   "fullname1",
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	rcd := httptest.NewRecorder()
	ui.WriteJsonResponse(status, data, rcd)

	//fmt.Println(rcd.Body.String())
	contentType := rcd.Header()["Content-Type"]
	s.Equal(contentType, []string{"application/json"})
}

func TestRunResp(t *testing.T) {
	suite.Run(t, new(_respSuite))
}
