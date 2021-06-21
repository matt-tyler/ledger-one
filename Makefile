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
