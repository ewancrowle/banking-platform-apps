module payment-switch-svc

go 1.26.2

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260209202127-80ab13bee0bf.1
	connectrpc.com/connect v1.19.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/moov-io/iso4217 v0.3.2
	google.golang.org/protobuf v1.36.11
)

require github.com/stretchr/testify v1.11.1 // indirect
