# GCP event driven rest API

This project demonstrates a "Hello, World!" example of triggering a Cloud Function via Pub/Sub when a REST API endpoint is called, focusing on local testing using emulators.


### Run PubSub emulator 

`$ gcloud beta emulators pubsub start`

* Note: This starts the emulator. You'll need to keep this running in a separate terminal window.

### Crete the topic 

`$ curl -X PUT http://localhost:8085/v1/projects/demo-test/topics/send-event-topic`

* Title: This command creates the Pub/Sub topic named "send-event-topic" in the emulator.
* Note: The demo-test project ID is arbitrary for the emulator.

### Check the Topic (Optional):

`$ curl http://localhost:8085/v1/projects/demo-test/topics`

* Title: This command lists the existing topics in the emulator, allowing you to verify that the topic was created successfully.


### Associate the topic with the Push endpoint:

```
curl -X PUT \
  http://localhost:8085/v1/projects/demo-test/subscriptions/my-subscription \
  -H "Content-Type: application/json" \
  -d '{
      "topic": "projects/demo-test/topics/send-event-topic",
      "pushConfig": {
        "pushEndpoint": "http://localhost:8082/projects/demo-test/topics/send-event-topic"
      }
  }'
```


### Check subscriptions

`$ curl http://localhost:8085/v1/projects/demo-test/subscriptions`

### Start the cloud event function

Navifgate to `/cmd/functions` and run

`$ go run main.go`

You should see the message:
2025/08/02 18:13:23 Listening on port 8082


### Call the REST API to trigger the entire flow:

`curl --location --request POST 'http://localhost:8080/save'`


### Payload to test the cloud event function
```
{
  "subscription": "projects/my-project/subscriptions/my-subscription",
  "message": {
    "@type": "type.googleapis.com/google.pubsub.v1.PubsubMessage",
    "attributes": {
      "attr1":"attr1-value"
    },
    "data": "dGVzdCBtZXNzYWdlIDM=",
    "messageId": "message-id",
    "publishTime":"2021-02-05T04:06:14.109Z"
  }
}

```

* the request to test the function
```
curl -X POST http://localhost:8082/projects/demo-test/topics/send-event-topic \
  -H "Content-Type: application/json" \
  -H "ce-id: test-123" \
  -H "ce-source: test-source" \
  -H "ce-specversion: 1.0" \
  -H "ce-type: type.googleapis.com/google.pubusb.v1.PubsubMessage" \
  -d '{
  "subscription": "projects/my-project/subscriptions/my-subscription",
  "message": {
    "@type": "type.googleapis.com/google.pubsub.v1.PubsubMessage",
    "attributes": {
      "attr1":"attr1-value"
    },
    "data": "eyJlbWFpbCI6ImNvbnRhY3RAeWFob28uY29tIn0=",
    "messageId": "message-id",
    "publishTime":"2021-02-05T04:06:14.109Z"
  }
}'
```