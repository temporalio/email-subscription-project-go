// @@@SNIPSTART subscription-workflow-go-activities
package subscribe_emails

import (
	"context"

	"go.temporal.io/sdk/activity"
)

type Activities struct {

}
// email activities
func (a *Activities) SendWelcomeEmail(ctx context.Context, emailInfo EmailInfo) (string, error) {
	activity.GetLogger(ctx).Info("sending welcome email to customer", emailInfo.EmailAddress)
	return "Sending welcome email completed for " + emailInfo.EmailAddress, nil
}

func (a *Activities) SendCancellationEmail(ctx context.Context, emailInfo EmailInfo) (string, error) {
	activity.GetLogger(ctx).Info("sending cancellation email during active subscription to: ", emailInfo.EmailAddress)
	return "Sending cancellation email during active subscription completed for: " + emailInfo.EmailAddress, nil
}

func (a *Activities) SendSubscriptionEndedEmail(ctx context.Context, emailInfo EmailInfo) (string, error) {
	activity.GetLogger(ctx).Info("sending subscription over email to: ", emailInfo.EmailAddress)
	return "Sending subscription over email completed for: " + emailInfo.EmailAddress, nil
}

func (a *Activities) SendSubscriptionEmail(ctx context.Context, emailInfo EmailInfo) (string, error) {
	activity.GetLogger(ctx).Info("sending subscription email to: ", emailInfo.EmailAddress)
	return "Sending subscription email for: " + emailInfo.EmailAddress, nil
}
// @@@SNIPEND