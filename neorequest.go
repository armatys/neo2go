package neo2go

import (
	"bytes"
	"net/http"
)

type NeoRequest struct {
	*http.Request
	batchId NeoBatchId
	result  interface{}
}

func NewNeoRequest(method, urlStr string, bodyData []byte, result interface{}) (*NeoRequest, error) {
	var (
		req *http.Request
		err error
	)
	if len(bodyData) > 0 {
		bodyBuf := bytes.NewBuffer(bodyData)
		req, err = http.NewRequest(method, urlStr, bodyBuf)
	} else {
		req, err = http.NewRequest(method, urlStr, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Stream", "true")
	if len(bodyData) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	neoRequest := new(NeoRequest)
	neoRequest.Request = req
	neoRequest.result = result
	return neoRequest, nil
}
