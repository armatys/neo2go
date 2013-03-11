package neo2go

import (
	"bytes"
	"net/http"
)

type NeoHttpRequest struct {
	*http.Request
}

func NewNeoHttpRequest(method, urlStr string, bodyBuf *bytes.Buffer) (*NeoHttpRequest, error) {
	var req *http.Request
	var err error

	if bodyBuf != nil {
		req, err = http.NewRequest(method, urlStr, bodyBuf)
	} else {
		req, err = http.NewRequest(method, urlStr, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Stream", "true")
	if bodyBuf != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	neoRequest := new(NeoHttpRequest)
	neoRequest.Request = req
	return neoRequest, nil
}
