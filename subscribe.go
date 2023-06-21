// @@@SNIPSTART subscription-workflow-go-subscribe
package subscribe_emails

var TaskQueueName string = "email_subscription"
var ClientHostPort string = "localhost:4000"

type EmailDetails struct {
	EmailAddress string
	Message string
	IsSubscribed bool
	MaxSubscriptionPeriods int
	SubscriptionCount int
}
// @@@SNIPEND