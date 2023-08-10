# subscription-workflow-go
The project files for the Subscription Workflow in Go tutorial.

# Temporal Subscription Workflow Tutorial - Go
Temporal promises to help you build long-lasting apps. With this example, you'll learn how to integrate key Temporal concepts into an email-based subscription program.

## Prerequisites
- Temporal Cloud
- Review Hello World

## Steps
Once this is working, follow the directions provided.

1. Launch an instance of Temporal Cloud (can be done on terminal or through a container such as Docker).
2. Enter `go run gateway/main.go` into a new terminal window.
3. Open another terminal and enter `go run worker/main.go`

## Curl commands

### subscribe

Use the curl command to send a POST request to `http://localhost:5000/subscribe` with the email address as a JSON payload.

```bash
curl -X POST -H "Content-Type: application/json" -d '{"email": "example@example.com"}' http://localhost:4000/subscribe
```

### get-details

The email address should be included in the query string parameter of the URL.

```bash
curl -X GET -H "Content-Type: application/json" http://localhost:4000/get_details?email=example@example.com

```

### Unsubscribe

Send a `DELETE` request with the email address in a JSON payload:

```bash
curl -X DELETE -H "Content-Type: application/json" -d '{"email": "example@example.com"}' http://localhost:4000/unsubscribe
```
