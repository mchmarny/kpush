GCP_PROJECT=s9-demo
GCP_REGION=us-central1
BINARY_NAME=pusheventing
PUBSUB_TOPIC=pusheventing
TOKEN=${KNWON_PUBLISHER_TOKEN}
HTTP_TARGET=https://msgme.default.knative.tech/push?publisherToken=${TOKEN}
DOCKER_USERNAME=mchmarny


topic:
	gcloud beta pubsub topics create ${PUBSUB_TOPIC}

subscription:
	gcloud beta pubsub subscriptions create ${PUBSUB_TOPIC}-sub \
    	--topic ${PUBSUB_TOPIC} \
    	--push-endpoint ${HTTP_TARGET} \
    	--ack-deadline 10

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

images: client-images server-images

image:
	gcloud builds submit \
		--project $(GCP_PROJECT) \
		--tag gcr.io/$(GCP_PROJECT)/${BINARY_NAME}-server:latest

docker:
	docker build -t $(BINARY_NAME) .
	docker tag $(BINARY_NAME):latest $(DOCKER_USERNAME)/$(BINARY_NAME):latest

docker-run:
	docker run -itP --expose 8080 $(DOCKER_USERNAME)/$(BINARY_NAME):latest

push:
	kubectl apply -f deploy/server.yaml