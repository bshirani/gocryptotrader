echo "GoCryptoTrader: Generating gRPC, proxy and swagger files."
# You may need to include the go mod package for the annotations file:
# $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.0.1/third_party/googleapis

export GOPATH=$(go env GOPATH)
export GAPI=$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis

protoc -I=. -I=$GOPATH/pkg/mod -I=$GAPI --go_out=. rpc.proto
protoc -I=. -I=$GOPATH/pkg/mod -I=$GAPI --go-grpc_out=. rpc.proto
protoc -I=. -I=$GOPATH/pkg/mod -I=$GAPI --grpc-gateway_out=logtostderr=true:. rpc.proto
protoc -I=. -I=$GOPATH/pkg/mod -I=$GAPI --openapiv2_out=logtostderr=true:. rpc.proto
