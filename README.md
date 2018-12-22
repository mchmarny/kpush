# pusheventing

> WIP: document not finished

Webhooks have experienced somewhat of a renaissance recently as more and more developers are learning about the benefits of serverless. Webhooks allow the developer to send "events" from one system to another based on triggers (e.g. GitHub fires Webhook on PR comment).

Because these events are most of the time sent across public network, Webhooks need to assure that the event was submitted from a valid source, and ideally, that the content of the event (the data) has not been tempered with in transit. This is where `pusheventing` comes in.

// TODO: diagram

This project demonstrates how to setup end-to-end pipeline using customer client to publish events onto [GCP Cloud PubSub](https://cloud.google.com/pubsub/) topic, and a service hosted on an instance of the [Knative](https://github.com/knative/) project to securely process these events

> NOTE: this is personal project for demonstration only.

## Setup

The set up the end-to-end pipeline we will have to configure:

* Service on Knative
* Topic and Subscription on PubSub
* Configure Client

### Service on Knative

// TODO: link to how to setup Knative

// TODO: instructions on how to deploy service

### Topic and Subscription on PubSub

// TODO: instructions on how to create topic

// TODO: instructions on how to create subscription

* https://cloud.google.com/pubsub/docs/subscriber

> Target of your PubSub push must be an HTTPS server with non-self-signed certificate

### Configure Client

// TODO: instructions on how to configure client and send data

```push
{
    "message": {
        "attributes": {
            "sig": "sha1=22c477fd1269c9d3bab8591b371a66976f10006e"
        },
        "data": "eyJpZC...",
        "messageId": "333651121184341",
        "publishTime": "2018-12-22T19:05:01.067Z",
    },
    "subscription": "projects/${PROJECT_ID}/subscriptions/pusheventing-push"
}
```



