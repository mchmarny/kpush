FROM golang:1.11.3
WORKDIR /go/src/github.com/mchmarny/pusheventing/cmd/server/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o app

FROM scratch
COPY --from=0 /go/src/github.com/mchmarny/pusheventing/cmd/server/app .
ENTRYPOINT ["/app"]