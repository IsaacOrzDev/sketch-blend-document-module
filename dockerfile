
FROM golang:1.19-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go get github.com/steebchen/prisma-client-go

RUN go run github.com/steebchen/prisma-client-go generate

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main .

EXPOSE 5003

CMD ["./main"]