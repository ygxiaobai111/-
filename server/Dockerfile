FROM golang:1.20-alpine AS builder

ENV GOPROXY https://goproxy.io,direct
ENV GO111MODULE on
WORKDIR $GOPATH/src/work
COPY . $GOPATH/src/work

RUN go mod tidy
RUN go install ./cmd/chatServer.go
RUN go install ./cmd/gameServer.go
RUN go install ./cmd/loginServer.go
RUN go install ./cmd/webServer.go
RUN go install ./cmd/gateServer.go


FROM alpine:3.18
WORKDIR /



COPY --from=builder /go/bin/chatServer  /bin/chatServer
COPY --from=builder /go/bin/gameServer  /bin/gameServer
COPY --from=builder /go/bin/loginServer  /bin/loginServer
COPY --from=builder /go/bin/webServer  /bin/webServer
COPY --from=builder /go/bin/gateServer  /bin/gateServer



COPY --from=builder   /go/src/work/config /config
EXPOSE 8004



CMD ["/bin/gateServer"]

