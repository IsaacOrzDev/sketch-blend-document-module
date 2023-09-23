gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/*.proto

db-pull:
	go run github.com/steebchen/prisma-client-go db pull

db-pull:
	go run github.com/steebchen/prisma-client-go db push	

db-generate:
	go run github.com/steebchen/prisma-client-go generate



