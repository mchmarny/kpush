# kpush - push signed messages from PubSub to Knative service

The use of [GCP Cloud PubSub](https://cloud.google.com/pubsub/) is common component in data flow pipelines. It creates global and elastic separation between the data provider and data consumers. One of its features is the ability to [push](https://cloud.google.com/pubsub/docs/push) subscription which sends messages to the configured `webhook`. Consider the following solution:

![kpush flow](img/kpush-flow.png)

In this case, the `kpush` client publishes messages to PubSub `topic` which the `subscription` subsequently pushes to target processing service, an app hosted on an instance of `Knative`. PubSub provides [instructions](https://cloud.google.com/pubsub/docs/push) on how to include [token](https://cloud.google.com/pubsub/docs/faq#security) in POST to ensure that only submissions from valid senders are processed by target service.

Since this token can leak (e.g. logs), `kpush` demonstrates how to add an additional level of validation by signing each message on the client and validating that signature on the service side. The below PubSub message sample illustrates the `sig` attribute holding [SHA-1](https://en.wikipedia.org/wiki/SHA-1) of the `message.data` payload which is appended to each message.

```json
{
    "message": {
        "attributes": {
            "sig": "sha1=22c477fd1269c9d3bab8591b371a66976f10006e"
        },
        "data": "eyJpZC...",
        "messageId": "333651121184341",
        "publishTime": "2018-12-22T19:05:01.067Z",
    }
}
```

## Setup

The `kpush` flow includes two components: `client`, generating, signing, and publishing mocked up messages, and `server` which receives PubSub pushed messages and validates their signature.

> Assuming `gcloud` SDK already configured. See [gcloud docs](https://cloud.google.com/sdk/gcloud/) for instructions if you need assistance

### PubSub topic

Let's start by creating the PubSub topic named `kpush` to which our client will be publishing messages.

```shell
gcloud pubsub topics create kpush
```

The response should be

```shell
Created topic [projects/YOUR_PROJECT_ID/topics/kpush].
```

### Knative Service

The installation and configuration of [Knative](https://github.com/knative) are beyond the scope of this readme, but, you can find detailed instructions on how to configure it on [Kubernetes](https://kubernetes.io/) service offered by most Cloud Service Providers [here](https://github.com/knative/docs/tree/master/install). In this example I will be using Google's [Kubernetes Engine](https://cloud.google.com/kubernetes-engine/) (GKE) which now can be easily configured with the validated version of Knative with a single checkbox.

Quickest way to build your service image is through [GCP Build](https://cloud.google.com/cloud-build/). Just submit the build request from within the `kpush` directory:

```shell
gcloud builds submit \
    --project ${GCP_PROJECT_ID} \
	--tag gcr.io/${GCP_PROJECT_ID}/kpush-server:latest
```

The build service is pretty verbose in output but eventually you should see something like this

```shell
ID           CREATE_TIME          DURATION  SOURCE                                   IMAGES                      STATUS
6905dd3a...  2018-12-23T03:48...  1M43S     gs://PROJECT_cloudbuild/source/15...tgz  gcr.io/PROJECT/kpush-server SUCCESS
```

The `IMAGE` column will have the repository URI for your new image (e.g. `gcr.io/PROJECT/kpush-server`)

Before we can deploy that service to Knative, we just need to update the `deploy/server.yaml` file

First, update the container `image` to the URI from the build `IMAGE` column value. Then update the two environment variables values:

* `KNOWN_PUBLISHER_TOKENS` which holds the token we will use in PubSub URL
* `MSG_SIG_KEY` which should be the key used by your clients to sign published messages

> Note, `KNOWN_PUBLISHER_TOKENS` may include many tokens separated by comma

Finally, to deploy `kpush` service to Knative apply the updated deployment using following command:

```shell
kubectl apply -f deploy/server.yaml
```

The response should be

```shell
service.serving.knative.dev "pushme" configured
```

To check if the service was deployed successfully you can check the status using `kubectl get pods` command. The response should look something like this (e.g. Ready `3/3` and Status `Running`).

```shell
NAME                                          READY     STATUS    RESTARTS   AGE
pushme-00002-deployment-5645f48b4d-mb24j      3/3       Running   0          4h
```

Knative uses convention to build serving URL by combining the deployment name (e.g. `pushme`), namespace name (e.g. `demo`), and the pre-configured domain name (e.g. `knative.tech`). The resulting URL should look something like this

```shell
https://pushme.default.knative.tech
```

Go ahead, test it in browser, you should following JSON response:

```json
{ "handlers": [ "POST: /post" ] }
```

Target of PubSub push must also be an HTTPS server with non-self-signed certificate. The instructions on how to configure Knative with SSL certificate are located [here](https://github.com/knative/docs/blob/master/serving/using-an-ssl-cert.md).

> Note, PubSub will publish only to registering endpoints which will require you to A) Verify you have access to the domain, and B) Register your domain in APIs & services. Instructions for both can be set up [here](https://cloud.google.com/pubsub/docs/push)

### PubSub Subscription

The final component we need to set up is the PubSub subscription which will push messages to our Knative hosted service. the `MYTOKEN` should be set to one of the tokens you defined in the Knative service manifest environment variables (`KNOWN_PUBLISHER_TOKENS`).

```shell
gcloud pubsub subscriptions create kpush-sub \
    --topic kpush \
    --push-endpoint https://pushme.demo.knative.tech/push?publisherToken=${MYTOKEN} \
    --ack-deadline 30
```

> Note, cold-starts on Knative can be sometimes slow so we added the `ack-deadline=30` parameter. For additional info on PubSub subscription options see this [doc](* https://cloud.google.com/pubsub/docs/subscriber)


## Run

To demo the entire pipeline, `kpush` includes a client which will mock up a few messages, sign them, and submit them to the configured topic.

Navigate to the `client` directory...

```shell
cd cmd/client
```

...and run the following command:

```shell
go run main.go --project ${GCP_PROJECT_ID} \
               --key ${MSG_SIG_KEY} \
               --topic kpush
```

The response should look something like this

```shell
2018/12/23 13:58:44 Using topic publisher: YOUR-GCP-PROJECT-NAME:kpush
   status: published msg[0] eacc9be0-06fd-11e9-868f-acde48001122
   status: published msg[1] eb2e0cae-06fd-11e9-868f-acde48001122
   status: published msg[2] eb39aaa0-06fd-11e9-868f-acde48001122
sent 3 messages
```

### Monitoring

A successful end-to-end run of the pipeline will result in `202` status code from the Knative hosted service. To validate you can either query the Knative logs or navigate to the [GCP Stackdriver](https://cloud.google.com/stackdriver/) console to review the PubSub subscription panel (Resources > PubSub > Subscriptions)

![kpush flow](img/kpush-chart.png)

## Direct Push

`kpush` also includes an `HTTPPublisher` which can skip the PubSub entirely and submit signed messages directly to the Knative service. To do that include `url` parameter in the `kpush client` command like this

```shell
go run main.go --project ${GCP_PROJECT_ID} \
               --key ${MSG_SIG_KEY} \
               --url https://pushme.demo.knative.tech/push?publisherToken=${MYTOKEN} \
               --messages 10
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

