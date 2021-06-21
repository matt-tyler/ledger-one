// +build tools

package tools

import (
	// golang code generator for protobuf
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	// the twirp service generator
	_ "github.com/twitchtv/twirp/protoc-gen-twirp"
)
