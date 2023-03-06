package subscribe_emails

import (
	"context"

	"go.temporal.io/sdk/activity"
)

type Activities struct {

}
// email activities
func (a *Activities) SendWelcomeEmail(ctx context.Context, subscription Subscription) (string, error) {
	activity.GetLogger(ctx).Info("sending welcome email to customer", subscription.EmailInfo.EmailAddress)
	return "Sending welcome email completed for " + subscription.EmailInfo.EmailAddress, nil
}

func (a *Activities) SendCancellationEmailDuringTrialPeriod(ctx context.Context, subscription Subscription) (string, error) {
	activity.GetLogger(ctx).Info("sending cancellation email during trial period to: ", subscription.EmailInfo.EmailAddress)
	return "Sending cancellation email during trial period completed for " + subscription.EmailInfo.EmailAddress, nil
}

func (a *Activities) ChargeCustomerForBillingPeriod(ctx context.Context, subscription Subscription) (string, error) {
	activity.GetLogger(ctx).Info("charging customer for billing period.")
	return "Charging for billing period completed for: " + subscription.EmailInfo.EmailAddress, nil
}

func (a *Activities) SendCancellationEmailDuringActiveSubscription(ctx context.Context, subscription Subscription) (string, error) {
	activity.GetLogger(ctx).Info("sending cancellation email during active subscription to: ", subscription.EmailInfo.EmailAddress)
	return "Sending cancellation email during active subscription completed for: " + subscription.EmailInfo.EmailAddress, nil
}

func (a *Activities) SendSubscriptionOverEmail(ctx context.Context, subscription Subscription) (string, error) {
	activity.GetLogger(ctx).Info("sending subscription over email to: ", subscription.EmailInfo.EmailAddress)
	return "Sending subscription over email completed for: " + subscription.EmailInfo.EmailAddress, nil
}

func (a *Activities) SendSubscriptionEmail(ctx context.Context, subscription Subscription) (string, error) {
	activity.GetLogger(ctx).Info("sending subscription email to: ", subscription.EmailInfo.EmailAddress)
	return "Sending subscription email completed for: " + subscription.EmailInfo.EmailAddress, nil
}
