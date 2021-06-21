package ledger

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	pb "github.com/matt-tyler/ledger-one/rpc/ledger"
)

func NewService(ddb dynamodb.Client) (*Server, error) {
	return &Server{ddb}, nil
}

type Server struct {
	ddb dynamodb.Client
}

func (s *Server) ClaimDomain(ctx context.Context, input *pb.ClaimDomainInput) (*pb.ClaimDomainOutput, error) {
	command := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: input.Root},
			"sk": &types.AttributeValueMemberS{Value: input.Subdomain},
		},
	}

	s.ddb.PutItem(ctx, &command)

	return &pb.ClaimDomainOutput{
		Domain: &pb.Domain{
			Root:      input.Root,
			Subdomain: input.Subdomain,
		},
	}, nil
}
