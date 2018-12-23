FROM golang:latest as build
ENV GO111MODULE=on

WORKDIR /go/src/github.com/mchmarny/kpush/
COPY . .

RUN go mod download

WORKDIR /go/src/github.com/mchmarny/kpush/cmd/server/

RUN CGO_ENABLED=0 go build -o /kpush

FROM scratch
COPY --from=build /kpush /app/
ENTRYPOINT ["/app/kpush"]