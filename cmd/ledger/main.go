// main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/smithy-go/middleware"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"

	"github.com/matt-tyler/ledger-one/internal/ledger"
	m "github.com/matt-tyler/ledger-one/internal/middleware"
	rpc "github.com/matt-tyler/ledger-one/rpc/ledger"
)

type APIGatewayProxyHandler = func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func createHandler(serveHTTP http.HandlerFunc) APIGatewayProxyHandler {
	handler := func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		reqAccessorV2 := core.RequestAccessor{}
		req, err := reqAccessorV2.ProxyEventToHTTPRequest(event)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
			}, nil
		}
		log.Println(req.Method, req.URL.String())
		writer := core.NewProxyResponseWriter()

		serveHTTP(writer, req.WithContext((ctx)))
		res, err := writer.GetProxyResponse()
		return res, err
	}
	return handler
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Panicf("Unable to load SDK config\n, %v", err)
	}

	tableName, ok := os.LookupEnv("TABLE_NAME")
	if !ok {
		log.Panicln("Failed to find required environment variable: TABLE_NAME")
	}

	withEndpoint := func(options *dynamodb.Options) {
		if endpoint, ok := os.LookupEnv("DDB_ENDPOINT"); ok {
			options.EndpointResolver = dynamodb.EndpointResolverFromURL(endpoint)
		}
	}

	ddb := dynamodb.NewFromConfig(cfg, withEndpoint, dynamodb.WithAPIOptions(func(stack *middleware.Stack) error {
		// Attach the custom middleware to the beginning of the Initialize step
		return stack.Initialize.Add(m.DefaultTableNameMiddleware(tableName), middleware.Before)
	}))

	l, _ := ledger.NewService(*ddb)
	service := rpc.NewLedgerServer(l)
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(createHandler(service.ServeHTTP))
}
