package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	muxadapter "github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
)

var muxLambda *muxadapter.GorillaMuxAdapter

func init() {
	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Cold start")

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello from /"))
	})

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello from /ping"))
	})

	router.HandleFunc("/ping/yes", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello from /ping/yes"))
	})

	router.HandleFunc("/serverless-env", func(w http.ResponseWriter, r *http.Request) {
		rfs := os.Getenv("RUNNING_FROM_SERVERLESS")
		if rfs != "true" {
			rfs = "false"
		}

		someEnv := os.Getenv("SOME_ENV")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{ "runningFromServerless": %s, "someEnv": "%s" }`, rfs, someEnv)))
	})

	muxLambda = muxadapter.New(router)
}

func Handler(ctx context.Context, event core.SwitchableAPIGatewayRequest) (*core.SwitchableAPIGatewayResponse, error) {
	return muxLambda.ProxyWithContext(ctx, event)
}

func main() {
	lambda.Start(Handler)
}
