package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"subscribe_emails"

	"go.temporal.io/sdk/client"
)

var temporalClient client.Client
var taskQueueName string

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "<h1>Sign up here!")
	_, _ = fmt.Fprint(w, "<form method='post' action='subscribe'><input required name='email' type='email'><input type='submit' value='Subscribe'>")
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
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

	// use the email as the id in the workflow. This may leak PII.
	workflowOptions := client.StartWorkflowOptions{
		ID:        "subscribe_email_" + email,
		TaskQueue: taskQueueName,
	}

	// Define the subscription
	subscription := subscribe_emails.Subscription{
		EmailInfo: subscribe_emails.EmailInfo{
			EmailAddress: email,
			Mail: "",
		},
		Periods: subscribe_emails.Periods{
			MaxBillingPeriods: 30,
			BillingPeriodCharge: 10,
		},
	}

	// execute the Temporal Workflow to start the subscription.
	_, err = temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, subscribe_emails.SubscriptionWorkflow, subscription)

	if err != nil {
		_, _ = fmt.Fprint(w, "<h1>Couldn't sign up</h1>")
		log.Print(err)
	} else {
		_, _ = fmt.Fprint(w, "<h1>Signed up!</h1>")
	}

}

func unsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	
	switch r.Method {

	case "GET":

		// http.ServeFile(w, r, "form.html")
		_, _ = fmt.Fprint(w, "<h1>Unsubscribe</h1><form method='post' action='/unsubscribe'><input required name='email' type='email'><input type='submit' value='Unsubscribe'>")

	case "POST":

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

		workflowID := "subscribe_email_" + email

		err = temporalClient.CancelWorkflow(context.Background(), workflowID, "")

		if err != nil {
			_, _ = fmt.Fprint(w, "<h1>Couldn't unsubscribe you</h1>")
			log.Fatalln("Unable to cancel Workflow Execution", err)
		} else {
			_, _ = fmt.Fprint(w, "<h1>Unsubscribed you from our emails. Sorry to see you go.</h1>")
			log.Println("Workflow Execution cancelled", "WorkflowID", workflowID)
		}
	}

}

func getDetailsHandler(w http.ResponseWriter, r *http.Request) {

	// http.ServeFile(w, r, "form.html")
	_, _ = fmt.Fprint(w, "<h1>Get description details here!</h1>")
	_, _ = fmt.Fprint(w, "<form method='get' action='/details'><input required name='email' type='email'><input type='submit' value='GetDetails'>")
	
}

func showDetailsHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "<h1>Your details:</h1>")
	_, _ = fmt.Fprint(w, "<p>"+"test"+"</p>")
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
	http.HandleFunc("/details", showDetailsHandler)
	_ = http.ListenAndServe(":"+port, nil)
}
