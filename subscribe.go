// @@@SNIPSTART subscription-workflow-go-subscribe
package subscribeemails

const TaskQueueName string = "email_subscription"
const ClientHostPort string = "localhost:4000"

type EmailDetails struct {
	EmailAddress      string `json:"emailAddress"`
	Message           string `json:"message"`
	IsSubscribed      bool   `json:"isSubscribed"`
	SubscriptionCount int    `json:"subscriptionCount"`
}

// @@@SNIPEND
