// @@@SNIPSTART subscription-workflow-go-worker
package main

import (
	"log"
	"subscribeemails"

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
	// create Worker
	w := worker.New(c, subscribeemails.TaskQueueName, worker.Options{})
	// register Activity and Workflow
	w.RegisterWorkflow(subscribeemails.SubscriptionWorkflow)
	w.RegisterActivity(subscribeemails.SendEmail)

	log.Println("Worker is starting.")
	// Listen to Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker.", err)
	}
}
// @@@SNIPEND