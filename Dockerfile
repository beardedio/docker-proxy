FROM golang as builder

ENV GO111MODULE=on

WORKDIR /code

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build


# final stage
FROM alpine
COPY --from=builder /code/docker-proxy /code/
EXPOSE 8080
ENTRYPOINT ["/code/docker-proxy"]
