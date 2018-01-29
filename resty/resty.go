package resty

import (
	"../utility"
	"gopkg.in/resty.v1"
)

type authSuccess struct {
	ID, Message string
}

type authError struct {
	ID, Message string
}

// HTTP Methods => Exported Type
const (
	GET     = iota
	POST    = iota
	PUT     = iota
	PATCH   = iota
	DELETE  = iota
	HEAD    = iota
	OPTIONS = iota
)

// Send => Exported
func Send(
	method int,
	url string,
	header map[string]string,
	body map[string]string,
	queryStrings map[string]string,
	token string) *resty.Response {

	req := resty.R().SetError(&authError{}).SetResult(&authSuccess{})

	if header != nil {
		req = req.SetHeaders(header)
	}

	if body != nil {
		req = req.SetBody(body)
	}

	if queryStrings != nil {
		req = req.SetQueryParams(queryStrings)
	}

	if token != "" {
		req = req.SetAuthToken(token)
	}

	var err error
	var resp *resty.Response
	resp = nil

	switch method {
	case GET:
		resp, err = req.Get(url)
		break

	case POST:
		resp, err = req.Post(url)
		break

	case PUT:
		resp, err = req.Put(url)
		break

	case PATCH:
		resp, err = req.Patch(url)
		break

	case DELETE:
		resp, err = req.Delete(url)
		break

	case HEAD:
		resp, err = req.Head(url)
		break

	case OPTIONS:
		resp, err = req.Options(url)
		break
	}

	if err != nil {
		utility.Log("(Resty=>Send) Failed to send request", err)
	}

	return resp
}
