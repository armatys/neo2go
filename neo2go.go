package neo2go

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type GraphDatabaseService struct {
	*NeoServiceRoot
	client  *http.Client
	selfURL string
}

func NewGraphDatabaseService(url string) *GraphDatabaseService {
	service := &GraphDatabaseService{
		NeoServiceRoot: new(NeoServiceRoot),
		client:         &http.Client{},
		selfURL:        url,
	}
	return service
}

func (self *GraphDatabaseService) Connect() (*NeoResponse, error) {
	req, err := NewNeoRequest("GET", self.selfURL, nil)
	if err != nil {
		return nil, err
	}

	neoResponse, err := self.execute(req, self.NeoServiceRoot)
	if err != nil {
		return nil, err
	}

	return neoResponse, nil
}

func (self *GraphDatabaseService) execute(neoRequest *NeoRequest, v interface{}) (*NeoResponse, error) {
	resp, err := self.client.Do(neoRequest.Request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var container interface{}

	neoResponse := new(NeoResponse)
	if resp.StatusCode >= 400 {
		container = neoResponse.NeoError
	} else {
		container = v
	}

	if container != nil {
		ctype := resp.Header.Get("content-type")
		if matched, err := regexp.MatchString("^application/json.*", ctype); matched {
			dec := json.NewDecoder(resp.Body)
			err = dec.Decode(container)
			if err != nil {
				return nil, err
			}
		} else if len(ctype) == 0 {
			return nil, fmt.Errorf("Server did not return a content-type for this response.")
		} else {
			return nil, fmt.Errorf("Server has returned a response with unsupported content-type (%s)", ctype)
		}
	}

	neoResponse.StatusCode = resp.StatusCode

	return neoResponse, nil
}
