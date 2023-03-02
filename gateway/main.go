package main

import (
	"fmt"
	"net/http"

	"go.temporal.io/sdk/client"
)

var temporalClient client.Client
var taskQueueName string

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "<h1>Sign up here!")
	_, _ = fmt.Fprint(w, "<form method='post' action='subscribe'><input required name='email' type='email'><input type='submit' value='Subscribe'>")
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	// compose email


	// define subscription


	// execute Workflow to start subsciption

}

func unsubscribeHandler(w http.ResponseWriter, r *http.Request) {


}

func getDetailsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Your details have been retrieved.")
	// create query
	// respond in email with results
}

func main() {
	port := "4000"
	taskQueueName = "subscription_emails"

	var err error
	temporalClient, err = client.Dial(client.Options {
		HostPort: client.DefaultHostPort,
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Starting the web server on port %s\n", port)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/subscribe", subscribeHandler)
	http.HandleFunc("/unsubscribe", unsubscribeHandler)
	http.HandleFunc("/getdetails", getDetailsHandler)

	e := http.ListenAndServe(":"+port, nil)
	if e != nil {
		http.ListenAndServe(":"+port, nil)
	}
}
