//go:build tools
// +build tools

//go:generate go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway

package tools

import _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
