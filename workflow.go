package subscribe_emails

import (
	"errors"
	"time"

	"go.temporal.io/sdk/workflow"
)

// Workflow definition
func SubscriptionWorkflow(ctx workflow.Context, subscription Subscription) error {

	billingPeriodNum := 0

// How frequently to send the messages
duration := 10 * time.Second

ao := workflow.ActivityOptions{
	StartToCloseTimeout: 10 * time.Minute,
	WaitForCancellation: true,
}

ctx = workflow.WithActivityOptions(ctx, ao)

var activities *Activities

logger := workflow.GetLogger(ctx)
logger.Info("Subscription created for " + subscription.EmailInfo.EmailAddress)

// Handle any cleanup, including cancellations.
defer func() {
	if !errors.Is(ctx.Err(), workflow.ErrCanceled) {
		return
	}

	// Cancellation received, which will trigger an unsubscribe email.

	newCtx, _ := workflow.NewDisconnectedContext(ctx)

	data := EmailInfo{
		EmailAddress: subscription.EmailInfo.EmailAddress,
		Mail: "You have been unsubscribed from the Subscription Workflow. Good bye.",
	}

	logger.Info("Sending unsubscribe email to " + subscription.EmailInfo.EmailAddress)
	err := workflow.ExecuteActivity(newCtx, activities.SendSubscriptionOverEmail, data).Get(newCtx, nil)

	if err != nil {
		logger.Error("Unable to send unsubscribe message", "Error", err)
	}
}()

logger.Info("Sending welcome email to " + subscription.EmailInfo.EmailAddress)

	data := EmailInfo{
		EmailAddress: subscription.EmailInfo.EmailAddress,
		Mail:         "Welcome! you've been signed up!",
	}

	err := workflow.ExecuteActivity(ctx, activities.SendWelcomeEmail, data).Get(ctx, nil)

	if err != nil {
		logger.Error("Failed to send welcome email", "Error", err)
	} else {
		billingPeriodNum++
	}

	// start subscription period
	for (billingPeriodNum < 10) {

		data := EmailInfo{
			EmailAddress: subscription.EmailInfo.EmailAddress,
			Mail:         "This is yet another email in the Subscription Workflow.",
		}

		err = workflow.ExecuteActivity(ctx, activities.SendSubscriptionEmail, data).Get(ctx, nil)

		if err != nil {
			logger.Error("Failed to send email ", "Error", err)
		}

		logger.Info("sent content email to " + subscription.EmailInfo.EmailAddress)

		billingPeriodNum++

		workflow.Sleep(ctx, duration)
	}

	return nil
}

