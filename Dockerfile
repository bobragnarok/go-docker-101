FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./cmd ./cmd

RUN go build -o /server ./cmd/main.go

EXPOSE 8080

CMD [ "/server" ]