// @@@SNIPSTART subscription-workflow-go-subscribe
package subscribeemails

const TaskQueueName string = "email_subscription"
const ClientHostPort string = "localhost:4000"

type EmailDetails struct {
	EmailAddress           string
	Message                string
	IsSubscribed           bool
	SubscriptionCount      int
}

// @@@SNIPEND
