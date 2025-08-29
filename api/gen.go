package api

//go:generate go tool oapi-codegen -o ../internal/controller/http/gen/types.go -generate types -package gen openapi.yaml
//go:generate go tool oapi-codegen  -o ../internal/controller/http/gen/server.go -generate chi-server,strict-server -package gen openapi.yaml
//go:generate go tool oapi-codegen  -o ../internal/controller/http/gen/spec.go -generate spec -package gen openapi.yaml

//go:generate go tool oapi-codegen -o ../pkg/client/types.go -generate types -package client openapi.yaml
//go:generate go tool oapi-codegen -o ../pkg/client/http_client.go -generate client -package client openapi.yaml
