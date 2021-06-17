package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/asim/go-micro/v3/errors"
	log "github.com/rs/zerolog/log"
)

func httpPost(endPoints []string, header map[string]string, body []byte) ([]byte, error) {
	if len(endPoints) == 0 {
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
		// certPath := "config/cert/localhost.crt"
		// keyPath := "config/cert/localhost.key"

		// tlsConf := &tls.Config{InsecureSkipVerify: true}
		// cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		// if err != nil {
		// 	log.Fatal().Err(err).Msg("load certificates failed")
		// }
		// tlsConf.Certificates = []tls.Certificate{cert}
		// http.DefaultClient.Transport = &http.Transport{TLSClientConfig: tlsConf}
		http.DefaultClient.Transport = &http.Transport{}
		http.DefaultClient.Timeout = time.Second * 3

		// make the call
		hrsp, err := http.DefaultClient.Do(hreq)
		if err != nil {
			return []byte(""), errors.InternalServerError(ep, err.Error())
		}

		// read body
		b, err := io.ReadAll(hrsp.Body)
		hrsp.Body.Close()
		if err != nil {
			return []byte(""), errors.InternalServerError(ep, err.Error())
		}

		if hrsp.StatusCode != 200 {
			return []byte(""), errors.InternalServerError(ep, fmt.Sprintf("%d", hrsp.StatusCode))
		}

		return b, nil
	}

	for _, endpoint := range endPoints {
		resp, err := fn(endpoint)
		if err != nil {
			log.Warn().Str("endpoint", endpoint).Err(err).Msg("http post failed")
			continue
		}

		return resp, err
	}

	return []byte(""), fmt.Errorf("no valid gate endpoint")
}
