package e2e

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	pb "github.com/matt-tyler/ledger-one/rpc/ledger"
	"google.golang.org/protobuf/testing/protocmp"
)

var urlFlag = flag.String("h", "http://localhost:8080", "url")

func TestMain(m *testing.M) {
	flag.Parse()
	fmt.Println("im here", *urlFlag)
	os.Exit(m.Run())
}

func TestApi(t *testing.T) {
	client := pb.NewLedgerProtobufClient(*urlFlag, &http.Client{})

	expected := &pb.ClaimDomainOutput{
		Domain: &pb.Domain{
			Root:      "example.cm",
			Subdomain: "myapp",
		},
	}

	result, err := client.ClaimDomain(context.Background(), &pb.ClaimDomainInput{
		Root:      "example.com",
		Subdomain: "myapp",
	})

	if err != nil {
		t.Error("Failed to create domain: ", err.Error())
	}

	diff := cmp.Diff(result.Domain, expected.Domain, protocmp.Transform())
	if diff != "" {
		t.Error("Unexpected result: ", diff)
	}
}
