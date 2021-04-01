package gate

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"testing"
	"time"

	"bitbucket.org/funplus/server/logger"
	"bitbucket.org/funplus/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	json "github.com/json-iterator/go"
	"github.com/urfave/cli/v2"
)

type header struct {
	Key   string
	Value string
}

var (
	ctx, _             = context.WithCancel(context.Background())
	ginServ *GinServer = nil
)

func newGinServer(t *testing.T) {
	if err := utils.RelocatePath("/server", "\\server"); err != nil {
		t.Fatalf("relocate path failed: %s", err.Error())
	}

	logger.InitLogger("gin_server_test")

	if ginServ != nil {
		return
	}

	set := flag.NewFlagSet("gin_test", flag.ContinueOnError)
	set.Bool("debug", true, "debug mode")
	set.String("cert_path_debug", "config/cert/localhost.crt", "cert path in debug mode")
	set.String("key_path_debug", "config/cert/localhost.key", "key path in debug mode")
	set.String("http_listen_addr", ":8080", "http listen address")
	set.String("https_listen_addr", ":4333", "https listen address")

	c := cli.NewContext(nil, set, nil)
	c.Context = ctx
	ginServ = NewGinServer(c, nil)

	ginServ.router.POST("/test_oneroute", func(ctx *gin.Context) {
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
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				log.Println(stack)
			}

			ginServ.Exit(c)
		}()

		if err := ginServ.Main(c); err != nil {
			log.Fatalf("gin server main failed: %s", err.Error())
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

	requestBody, _ := json.Marshal(map[string]string{
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
	newGinServer(t)

	time.Sleep(time.Second)
	performRequest(t, "POST", "http://localhost:8080/test_oneroute")

	result := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			performRequest(t, "POST", "http://localhost:8080/test_oneroute")
		}
	})

	fmt.Println("gin server benchmark result: ", result.String(), result.MemString())
}
