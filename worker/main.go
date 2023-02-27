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

	w := worker.New(c, "SUBSCRIPTION_TASK_QUEUE", worker.Options{})
	// register Activity and Workflow
	w.RegisterActivity(subscribe_emails.SendContentEmail)
	w.RegisterWorkflow(subscribe_emails.SubscriptionWorkflow)

	// Listen to Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker.", err)
	}
	log.Println("Worker successfully started.")
}