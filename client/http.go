package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/micro/go-micro/errors"
	logger "github.com/sirupsen/logrus"
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

		// cert
		certPath := "config/cert/localhost.crt"
		keyPath := "config/cert/localhost.key"

		tlsConf := &tls.Config{InsecureSkipVerify: true}
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			logger.Fatal("load certificates failed:", err)
		}
		tlsConf.Certificates = []tls.Certificate{cert}
		http.DefaultClient.Transport = &http.Transport{TLSClientConfig: tlsConf}

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
		resp, err := fn(endpoint)
		if err != nil {
			logger.WithFields(logger.Fields{
				"endpoint": endpoint,
				"error":    err,
			}).Warn("http post failed")
			continue
		}

		return resp, err
	}

	return []byte(""), fmt.Errorf("no valid gate endpoint")
}
