# pusheventing

WIP: Message signing demo with GCP PubSub client and Knative service

> Target of your PubSub push must be an HTTPS server with non-self-signed certificate


## Message

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

## Links

* https://cloud.google.com/pubsub/docs/subscriber
