package ui

import (
	"encoding/json"
	"log"
	"net/http"
)

type Status int

const (
	StatusOK Status = iota
	StatusUserExisted
	StatusInvalidContent
)

func (status Status) String() string {
	switch status {
	case StatusOK:
		return "Success"
	case StatusUserExisted:
		return "The user to be signed up has been existed"
	case StatusInvalidContent:
		return "The content is invalid"
	default:
		return ""
	}
}

type Info struct {
	Status  `json:"status"`
	Message string `json:"message"`
}

type Response struct {
	Info `json:"info"`
	Data interface{} `json:"data"`
}

func (resp *Response) JsonIdent() ([]byte, error) {
	return json.MarshalIndent(resp, "", "    ")
}

func WriteJsonResponse(status Status, data interface{}, w http.ResponseWriter) {
	switch status {
	case StatusUserExisted:
		w.WriteHeader(http.StatusNotAcceptable)
	case StatusInvalidContent:
		w.WriteHeader(http.StatusBadRequest)
	}

	resp := Response{
		Info: Info{
			Status:  status,
			Message: status.String(),
		},
		Data: data,
	}

	body, err := resp.JsonIdent()
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(body); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header()["Content-Type"] = []string{"application/json"}
}
