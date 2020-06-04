module github.com/enixdark/sample

go 1.14

require (
	github.com/golang/protobuf v1.4.2
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v2 v2.8.0
	github.com/micro/micro/v2 v2.8.0 // indirect
	github.com/micro/protoc-gen-micro/v2 v2.3.0 // indirect
	go.mongodb.org/mongo-driver v1.3.4 // indirect
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9 // indirect
	golang.org/x/sys v0.0.0-20200602225109-6fdc65e7d980 // indirect
	google.golang.org/genproto v0.0.0-20200603110839-e855014d5736 // indirect
	google.golang.org/grpc v1.29.1
)

replace google.golang.org/grpc v1.29.1 => google.golang.org/grpc v1.26.0
