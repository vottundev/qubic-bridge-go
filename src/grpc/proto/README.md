https://pascalallen.medium.com/how-to-build-a-grpc-server-in-go-943f337c4e05

protoc -I=./ --go_out=./ --go_opt=Morders.proto=./ orders.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative orders.proto