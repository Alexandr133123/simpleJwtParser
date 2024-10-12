FROM golang:1.22

WORKDIR /app

COPY go.mod ./
RUN go mod download && go mod verify

EXPOSE 8080

COPY . .
RUN go build -v -o /jwtParser ./...

ENTRYPOINT ["/jwtParser"]