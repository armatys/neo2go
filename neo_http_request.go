package neo2go

import (
	"bytes"
	"fmt"
	"net/http"
)

var versionHeader string

func init() {
	versionHeader = fmt.Sprintf("neo2go v.%d", Version)
}

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
	req.Header.Set("User-Agent", "neo2go v.%d")
	if bodyBuf != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	neoRequest := new(NeoHttpRequest)
	neoRequest.Request = req
	return neoRequest, nil
}
