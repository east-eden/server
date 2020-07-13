package gate

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/urfave/cli/v2"
)

type header struct {
	Key   string
	Value string
}

var (
	ctx, cancel            = context.WithCancel(context.Background())
	ginServ     *GinServer = nil
)

func init() {
	os.Chdir("../")
}

func newGinServer() {
	if ginServ != nil {
		return
	}

	set := flag.NewFlagSet("benchmark_test", flag.ContinueOnError)
	set.Bool("debug", true, "debug mode")
	set.String("cert_path_debug", "config/cert/localhost.crt", "cert path in debug mode")
	set.String("key_path_debug", "config/cert/localhost.key", "key path in debug mode")
	set.String("http_listen_addr", ":8080", "http listen address")
	set.String("https_listen_addr", ":4433", "https listen address")

	c := cli.NewContext(nil, set, nil)
	c.Context = ctx
	ginServ = NewGinServer(nil, c)

	ginServ.engine.POST("/test_oneroute", func(ctx *gin.Context) {
		var req struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}

		if ctx.Bind(&req) == nil {
			diff := cmp.Diff(req.Value, "value_1001")
			if diff == "" {
				ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
			}
			return
		}

		ctx.String(http.StatusBadRequest, "bad request")
	})

	go func() {
		defer ginServ.Exit(c)

		if err := ginServ.Main(c); err != nil {
			log.Fatal(err)
		}
	}()
}

func performRequest(t *testing.T, method, url string, headers ...header) {
	tr := &http.Transport{
		//TLSClientConfig: &tls.Config{
		//InsecureSkipVerify: true,
		//},
	}
	client := &http.Client{Transport: tr}

	requestBody, err := json.Marshal(map[string]string{
		"key":   "1001",
		"value": "value_1001",
	})

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("performRequest failed: %s", err.Error())
	}
	defer resp.Body.Close()

	diff := cmp.Diff(resp.StatusCode, http.StatusOK)
	if diff != "" {
		t.Errorf("performRequest failed: %s", diff)
	}
}

func TestOneRoute(t *testing.T) {
	newGinServer()

	time.Sleep(time.Second)
	performRequest(t, "POST", "http://localhost:8080/test_oneroute")

	testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			performRequest(t, "POST", "http://localhost:8080/test_oneroute")
		}
	})
}
