package subscribe_emails

// EmailInfo is the data that the SendContentEmail uses to send the message.
type EmailInfo struct {
	EmailAddress string
	Mail         string
}

// Campaign is the info about the email campaign.
type Campaign struct {
	Name             string
	WelcomeEmail     string
	UnsubscribeEmail string
	Mails            []string
}

// Subscription is the user email and the campaign they'll receive.
type Subscription struct {
	EmailAddress string
	Campaign     Campaign
}
