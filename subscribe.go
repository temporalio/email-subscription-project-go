// @@@SNIPSTART subscription-workflow-go-subscribe
package subscribe_emails

import "time"

// EmailInfo is the data that the SendContentEmail uses to send the message.
type EmailInfo struct {
	EmailAddress string
	Mail string
}

// Periods contains duration info for trial and billing periods
type Periods struct {
	SubcriptionPeriod time.Duration
	MaxSubscriptionPeriods int
}

// Subscription is the user email and the campaign they'll receive.
type Subscription struct {
	EmailInfo    EmailInfo
	Periods      Periods
}
// @@@SNIPEND