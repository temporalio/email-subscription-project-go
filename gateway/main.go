// @@@SNIPSTART subscription-workflow-go-gateway
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"subscribeemails"

	"go.temporal.io/sdk/client"
)

var temporalClient client.Client

// create the index handler, accessed at localhost:4000
func indexHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "<h1>Sign up here!")
	_, _ = fmt.Fprint(w, "<form method='post' action='/subscribe'><input required name='email' type='email'><input type='submit' value='Subscribe'>")
}

// create subscribe handler, which collects the email in the index handler form
func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// in case of any error
		_, _ = fmt.Fprint(w, "<h1>Error processing form</h1>")
		return
	}
	// check for valid email value
	email := r.PostForm.Get("email")

	if email == "" {
		// in case of any error
		_, _ = fmt.Fprint(w, "<h1>Email is blank</h1>")
		return
	}

	// use the email as the id in the workflow.
	workflowOptions := client.StartWorkflowOptions{
		ID:        email,
		TaskQueue: subscribeemails.TaskQueueName,
		WorkflowExecutionErrorWhenAlreadyStarted: true,
	}

	// Define the EmailDetails struct
	subscription := subscribeemails.EmailDetails {
		EmailAddress: email,
		Message: "Welcome to the Subscription Workflow!",
		SubscriptionCount: 0,
		MaxSubscriptionPeriods: 12,
	}

	// Execute the Temporal Workflow to start the subscription.
	_, err = temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, subscribeemails.SubscriptionWorkflow, subscription)

	if err != nil {
		_, _ = fmt.Fprint(w, "<h1>Couldn't sign up user. Please try again.</h1>")
		log.Print(err)
	} else {
		_, _ = fmt.Fprint(w, "<h1>Signed up! Resource was created successfully.</h1>")
	}

}

// create unsubscribe handler, accessed at localhost:4000/unsubscribe
func unsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	
	switch r.Method {

	case "GET":
		// create input field for the email
		_, _ = fmt.Fprint(w, "<h1>Unsubscribe</h1><form method='post' action='/unsubscribe'><input required name='email' type='email'><input type='submit' value='Unsubscribe'>")

	case "POST":
		// check value in input field
		err := r.ParseForm()

		if err != nil {
			// in case of any error
			_, _ = fmt.Fprint(w, "<h1>Error processing form</h1>")
			return
		}

		email := r.PostForm.Get("email")

		if email == "" {
			// in case of any error
			_, _ = fmt.Fprint(w, "<h1>Email is blank</h1>")
			return
		}
		// get Workflow ID to unsubscribe
		workflowID := email
		// cancel the Workflow Execution
		err = temporalClient.CancelWorkflow(context.Background(), workflowID, "")
		if err != nil {
			_, _ = fmt.Fprint(w, "<h1>Couldn't unsubscribe the user.</h1>")
			log.Fatalln("Unable to cancel Workflow Execution", err)
		} else {
			_, _ = fmt.Fprint(w, "<h1>Unsubscribed you from our emails. Sorry to see you go.</h1>")
			log.Println("Workflow Execution cancelled", "WorkflowID", workflowID)
		}
	}

}

// create part of the Query handler, accessed at localhost:4000/getdetails
func getDetailsHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "<h1>Get subscription details here!</h1>")
	_, _ = fmt.Fprint(w, "<form method='get' action='/details'><input required name='email' type='email'><input type='submit' value='GetDetails'>")
}

// create part of the Query handler that returns information at localhost:4000/details
func showDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the query string
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
	 	log.Println("Failed to query Workflow.")
		return 
	}
   
	// Extract the email parameter
	email := queryValues.Get("email")

	workflowID := email
	queryType := "GetDetails"
	
	// print email, billing period, charge, etc.
	resp, err := temporalClient.QueryWorkflow(context.Background(), workflowID, "", queryType)
	if err != nil {
		log.Fatalln("Unable to Query Workflow", err)
	}
	var result string
	if err := resp.Get(&result); err != nil {
		log.Fatalln("Unable to decode Query result", err)
	}
	log.Println("Received Query result", "Result: " + result)
	fmt.Fprint(w, "Your details have been retrieved. Results: " + result)
}
// set up handlers, Client, port, Task Queue name.
func main() {

	var err error
	temporalClient, err = client.Dial(client.Options {
		HostPort: client.DefaultHostPort,
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Starting the web server on %s\n", subscribeemails.ClientHostPort)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/subscribe", subscribeHandler)
	http.HandleFunc("/unsubscribe", unsubscribeHandler)
	http.HandleFunc("/getdetails", getDetailsHandler)
	http.HandleFunc("/details", showDetailsHandler)
	_ = http.ListenAndServe(":4000", nil)
}
// @@@SNIPEND