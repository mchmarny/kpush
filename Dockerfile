FROM golang:latest as build
ENV GO111MODULE=on

WORKDIR /go/src/github.com/mchmarny/pusheventing/
COPY . .

RUN go mod tidy

WORKDIR /go/src/github.com/mchmarny/pusheventing/cmd/server/

RUN CGO_ENABLED=0 go build

FROM scratch
COPY --from=build \
    /go/src/github.com/mchmarny/pusheventing/cmd/server/server \
    /app/
ENTRYPOINT ["/app/server"]