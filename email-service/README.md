# Email Service - Testing Kafka Message Injection

## Sending Kafka Message with kcat (Recommended)

If you have [kcat](https://github.com/edenhill/kcat) (formerly kafkacat) installed, you can send a JSON message directly
from file:

```sh
kcat -b localhost:9092 -t notification_email_topic -P notification-service/email-service/mock-payload.json
```

Or using Docker, if you don't have kcat on your host:

```sh
docker run --rm -v "$PWD":/data --network="host" edenhill/kcat:1.7.0 \
  -b localhost:9092 -t notification_email_topic -P /data/notification-service/email-service/mock-payload.json
```

## Sending Message via HTTP Endpoint (if microservice exposes producer)

If you add a POST /produce endpoint to one of your microservices, you can use curl:

```sh
curl -X POST -H "Content-Type: application/json" \
     --data-binary @notification-service/email-service/mock-payload.json \
     http://localhost:<your-producer-port>/produce
```

(**Note:** You must implement a /produce endpoint for this to work!)

## Other Methods

You can also use the Kafka console producer (see main project README), or write a short producer script in Python, Go,
or Node.

---
**mock-payload.json** contains a valid payload to test your email topic integration.
