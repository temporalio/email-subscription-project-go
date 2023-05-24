// @@@SNIPSTART subscription-workflow-go-subscribe
package subscribe_emails

import "time"

// EmailInfo is the data that the SendContentEmail uses to send the message.
type EmailInfo struct {
	EmailAddress string
	Mail string
}

// Subscription is the user email and duration information.
type Subscription struct {
	EmailInfo    EmailInfo
	SubscriptionPeriod time.Duration
	MaxSubscriptionPeriods int
}
// @@@SNIPEND