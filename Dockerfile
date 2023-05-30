FROM golang:1.19

# TODO builder pattern
RUN mkdir /kubem
WORKDIR /kubem

COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o /build

EXPOSE 9000
RUN pwd
ENTRYPOINT go run main.go
