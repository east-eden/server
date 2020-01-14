package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/micro/go-micro/errors"
)

func httpPost(c *TcpClient, header map[string]string, body []byte) ([]byte, error) {
	if len(c.gateEndpoints) == 0 {
		return []byte(""), fmt.Errorf("gate endpoints empty")
	}

	fn := func(ep string) ([]byte, error) {
		// send to backend
		hreq, err := http.NewRequest("POST", ep, bytes.NewReader(body))
		if err != nil {
			return []byte(""), errors.InternalServerError(ep, err.Error())
		}

		for k, v := range header {
			hreq.Header.Set(k, v)
		}

		// make the call
		hrsp, err := http.DefaultClient.Do(hreq)
		if err != nil {
			return []byte(""), errors.InternalServerError(ep, err.Error())
		}

		// read body
		b, err := ioutil.ReadAll(hrsp.Body)
		hrsp.Body.Close()
		if err != nil {
			return []byte(""), errors.InternalServerError(ep, err.Error())
		}

		if hrsp.StatusCode != 200 {
			return []byte(""), errors.InternalServerError(ep, string(hrsp.StatusCode))
		}

		return b, nil
	}

	for _, endpoint := range c.gateEndpoints {
		if resp, err := fn(endpoint); err == nil {
			return resp, err
		}
	}

	return []byte(""), fmt.Errorf("no valid gate endpoint")
}
