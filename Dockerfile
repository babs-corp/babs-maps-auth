FROM golang:1.23-alpine3.20

WORKDIR /src

COPY . .

RUN GO111MODULE="on" CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o app ./cmd/sso

EXPOSE 8082

ENV CONFIG_PATH=./config/dev.yaml  

CMD ["./app"]
