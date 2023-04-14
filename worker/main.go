// @@@SNIPSTART subscription-workflow-go-worker
package main

import (
	"log"
	"subscribe_emails"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// create client and worker
	c, err := client.Dial(client.Options {
		HostPort: client.DefaultHostPort,
		Namespace: client.DefaultNamespace,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal Client.", err)
	}
	defer c.Close()

	w := worker.New(c, "subscription_emails", worker.Options{})
	// register Activity and Workflow
	w.RegisterWorkflow(subscribe_emails.SubscriptionWorkflow)
	w.RegisterActivity(&subscribe_emails.Activities{})

	// Listen to Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker.", err)
	}
	log.Println("Worker successfully started.")
}
// @@@SNIPEND