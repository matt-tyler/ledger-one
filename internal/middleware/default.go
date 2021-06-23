package middleware

import (
	"context"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/smithy-go/middleware"
)

func reflectStructField(Iface interface{}, FieldName string) (*reflect.Value, error) {
	ValueIface := reflect.ValueOf(Iface)

	// Check if the passed interface is a pointer
	if ValueIface.Type().Kind() != reflect.Ptr {
		// Create a new type of Iface's Type, so we have a pointer to work with
		ValueIface = reflect.New(reflect.TypeOf(Iface))
	}

	// 'dereference' with Elem() and get the field by name
	Field := ValueIface.Elem().FieldByName(FieldName)
	if !Field.IsValid() {
		return nil, fmt.Errorf("Interface `%s` does not have the field `%s`", ValueIface.Type(), FieldName)
	}
	return &Field, nil
}

func DefaultTableNameMiddleware(tableName string) middleware.InitializeMiddleware {
	middleware := middleware.InitializeMiddlewareFunc("defaultTableName", func(
		ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
	) (
		out middleware.InitializeOutput, metadata middleware.Metadata, err error,
	) {
		if field, err := reflectStructField(in.Parameters, "TableName"); err == nil {
			field.Set(reflect.ValueOf(aws.String(tableName)))
		}

		return next.HandleInitialize(ctx, in)
	})
	return middleware
}
