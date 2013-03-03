package neo2go

import (
	"io"
	"net/http"
)

type NeoRequest struct {
	*http.Request
}

func NewNeoRequest(method, urlStr string, body io.Reader) (*NeoRequest, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Stream", "true")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	neoRequest := new(NeoRequest)
	neoRequest.Request = req
	return neoRequest, nil
}
