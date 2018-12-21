GCP_PROJECT=s9-demo
GCP_REGION=us-central1
PUBSUB_TOPIC=pusheventing

topic:
	gcloud beta pubsub topics create ${PUBSUB_EVENTS_TOPIC}


test:
	go test ./... -v

cover:
	go test ./... -cover
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

deps:
	go mod tidy

docs:
	godoc -http=:8888 &
	open http://localhost:8888/pkg/github.com/mchmarny/pusheventing/
	# killall -9 godoc