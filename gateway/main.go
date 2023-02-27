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
	// compose email
	err := r.ParseForm()
	if err != nil {
		_, _ = fmt.Fprint(w, "<h1>Error processing form</h1>")
		return
	}

	email := r.PostForm.Get("email")
	
	if email == "" {
		_, _ = fmt.Fprint(w, "<h1>Email is blank</h1>")
		return
	}
	
	workflowOptions := client.StartWorkflowOptions {
		ID: "email_drip_" + email,
		TaskQueue: taskQueueName,
	}

	// define subscription
	subscription := subscribe_emails.Subscription {
		EmailAddress: email,
		Campaign: subscribe_emails.Campaign{
			Name: "Subscription email",
			WelcomeEmail: "../mails/welcome.md",
			UnsubscribeEmail: "../mails/goodbye.md",
			Mails: []string{"../mails/1.md", "../mails/2.md", "../mails/3.md"},
		},

	}

	// execute Workflow to start subsciption
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
		_, _ = fmt.Fprint(w, "<h1> Unsubscribe</h1><form method='post' action='unsubscribe'><input required name='email' type='email><input type='submit' value='Unsubscrube'>")
	
	case "POST":
		err := r.ParseForm()

		if err != nil {
			_, _ = fmt.Fprint(w, "<h1>Error processing form</h1>")
			return
		}

		email := r.PostForm.Get("email")

		if email == "" {
			_, _ = fmt.Fprint(w, "<h1>Email is blank</h1>")
			return
		}

		workflowID := "email_drip_" + email

		err = temporalClient.CancelWorkflow(context.Background(), workflowID, "")

		if err != nil {
			_, _ = fmt.Fprint(w, "<h1>Couldn't unsubscribe the user</h1>")
			log.Fatalln("Unable to cancel Workflow Execution", err)
		} else {
			_, _ = fmt.Fprint(w, "<h1>Unsubscribed you from our emails. Sorry to see you go.</h1>")
			log.Println("WorkflowExecution cancelled", "WorkflowID", workflowID)
		}
	}
}

func getDetailsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Your details have been retrieved.")
	// create query
	// respond in email with results
}

func main() {
	port := "7233"
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

	e := http.ListenAndServe(":7233", nil)
	if e != nil {
		http.ListenAndServe(":"+port, nil)
	}
}
