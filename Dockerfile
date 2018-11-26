FROM golang:latest AS build
WORKDIR /go/src/github.com/viktorminko/monitor
ADD ./pkg pkg/
ADD ./cmd cmd/
RUN cd cmd && go get && CGO_ENABLED=0 go build -o /monitor

FROM alpine
RUN apk update && apk add ca-certificates
WORKDIR /app
COPY --from=build /monitor monitor
CMD ["/app/monitor", "-workdir", "/app/config"]