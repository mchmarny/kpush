GCP_PROJECT=s9-demo
BINARY_NAME=kpush
PUBSUB_TOPIC=$(BINARY_NAME)
TOKENS=${KNWON_PUBLISHER_TOKEN}
HTTP_TARGET=https://pushme.default.knative.tech/push?publisherToken=${PUBLISHER_TOKEN}


topic:
	gcloud beta pubsub topics create $(PUBSUB_TOPIC)

push:
	gcloud beta pubsub subscriptions create $(PUBSUB_TOPIC)-sub \
    	--topic $(PUBSUB_TOPIC) \
    	--push-endpoint $(HTTP_TARGET) \
    	--ack-deadline 30

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
	open http://localhost:8888/pkg/github.com/mchmarny/$(BINARY_NAME)/
	# killall -9 godoc

image:
	gcloud builds submit \
		--project $(GCP_PROJECT) \
		--tag gcr.io/$(GCP_PROJECT)/$(BINARY_NAME)-server:latest

docker:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker run -itP --expose 8080 $(DOCKER_USERNAME)/$(BINARY_NAME):latest

service:
	kubectl apply -f deploy/server.yaml