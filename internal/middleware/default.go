package middleware

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	sdk "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/smithy-go/middleware"
)

func DefaultTableNameMiddleware(tableName string) middleware.InitializeMiddleware {
	middleware := middleware.InitializeMiddlewareFunc("defaultTableName", func(
		ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
	) (
		out middleware.InitializeOutput, metadata middleware.Metadata, err error,
	) {
		if sdk.GetServiceID(ctx) != dynamodb.ServiceID {
			return next.HandleInitialize(ctx, in)
		}
		switch v := in.Parameters.(type) {
		case *dynamodb.PutItemInput:
			if v.TableName == nil {
				v.TableName = aws.String(tableName)
			}
		}

		return next.HandleInitialize(ctx, in)
	})
	return middleware
}
