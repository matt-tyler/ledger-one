
# These environment variables must be set for deployment to work.
S3_BUCKET := $(S3_BUCKET)
STACK_NAME := $(STACK_NAME)

# Common values used throughout the Makefile, not intended to be configured.
TEMPLATE = template.yaml
PACKAGED_TEMPLATE = packaged.yaml

download:
	@echo Download go.mod dependencies
	@go mod download

install-tools: download
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

gen:
	# Auto-generate code
	protoc \
		-I . \
		--twirp_out=. \
		--go_out=. rpc/ledger/service.proto

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: build
build:
	go build -o bin/ledger cmd/ledger/main.go

.PHONY: build-lambda
build-lambda:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-a -installsuffix cgo \
		-o bin/lambda/ledger cmd/ledger/main.go

.PHONY: package
package: build
	sam package \
		--template-file $(TEMPLATE) \
		--s3-bucket $(S3_BUCKET) \
		--output-template-file $(PACKAGED_TEMPLATE)

.PHONY: deploy
deploy: package
	sam deploy \
		--stack-name $(STACK_NAME) \
		--template-file $(PACKAGED_TEMPLATE) \
		--capabilities CAPABILITY_IAM

.PHONY: teardown
teardown:
	aws cloudformation delete-stack --stack-name $(STACK_NAME)