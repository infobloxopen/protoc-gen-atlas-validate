FROM golang:1.10.0 AS builder

LABEL stage=server-intermediate

WORKDIR /go/src/github.com/infobloxopen/protoc-gen-atlas-validate
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/usr/bin/protoc-gen-atlas-validate main.go

FROM infoblox/atlas-gentool:latest AS runner

COPY --from=builder /out/usr/bin/protoc-gen-atlas-validate /usr/bin/protoc-gen-atlas-validate
COPY --from=builder /go/src/github.com/infobloxopen/protoc-gen-atlas-validate/options/*.proto /go/src/github.com/infobloxopen/protoc-gen-atlas-validate/options/
COPY --from=builder /go/src/github.com/infobloxopen/protoc-gen-atlas-validate/example/examplepb/*.proto /go/src/github.com/infobloxopen/protoc-gen-atlas-validate/example/examplepb/
COPY --from=builder /go/src/github.com/infobloxopen/protoc-gen-atlas-validate/example/external/*.proto /go/src/github.com/infobloxopen/protoc-gen-atlas-validate/example/external/

WORKDIR /go/src
