// @@@SNIPSTART subscription-workflow-go-workflow
package subscribe_emails

import (
	"errors"
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// Workflow definition
func SubscriptionWorkflow(ctx workflow.Context, subscription Subscription) error {
	// declare variables, duration, and logger
	var activities *Activities
	subscriptionPeriod := 0
	// duration can be set up to a month.
	duration := 10 * 24 * time.Hour

	logger := workflow.GetLogger(ctx)
	logger.Info("Subscription created", "EmailAddress", subscription.EmailInfo.EmailAddress)
	// Query handler
	err := workflow.SetQueryHandler(ctx, "GetDetails", func(input []byte) (string, error) {
		return fmt.Sprintf("%v is on billing period %v out of %v",
			subscription.EmailInfo.EmailAddress,
			subscriptionPeriod,
			subscription.MaxSubscriptionPeriods), nil
	})
	if err != nil {
		return err
	}
	// variable for Activity Options. Timeout can be set to a longer timespan (such as a month)
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
		WaitForCancellation: true,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval: 2 * time.Hour,
			MaximumAttempts: 5,
		},
	})

	// Handle any cleanup, including cancellations.
	defer func() {
		newCtx, cancel := workflow.NewDisconnectedContext(ctx)
		defer cancel()

		if errors.Is(ctx.Err(), workflow.ErrCanceled) {
			data := EmailInfo{
				EmailAddress: subscription.EmailInfo.EmailAddress,
				Mail:         "Oh my! Looks like your subscription has been canceled!",
			}
			// send cancellation email
			err := workflow.ExecuteActivity(newCtx, activities.SendCancellationEmail, data)
			if err != nil {
				logger.Error("Failed to send cancel email", "Error", err)
			} else {
				// Cancellation received, which will trigger an unsubscribe email.
				logger.Info("Sent cancellation email")
			}
		}

		// information for the newly-ended subscription email
		data := EmailInfo{
			EmailAddress: subscription.EmailInfo.EmailAddress,
			Mail:         "You have been unsubscribed from the Subscription Workflow. Good bye.",
		}
		logger.Info("Sending unsubscribe email", "EmailAddress", subscription.EmailInfo.EmailAddress)
		// send the cancelled subscription email
		err := workflow.ExecuteActivity(newCtx, activities.SendSubscriptionEndedEmail, data).Get(newCtx, nil)

		if err != nil {
			logger.Error("Unable to send unsubscribe message", "Error", err)
		}
	}()
	// handling for the first email ever
	logger.Info("Sending welcome email", "EmailAddress", subscription.EmailInfo.EmailAddress)

	data := EmailInfo{
		EmailAddress: subscription.EmailInfo.EmailAddress,
		Mail:         "Welcome! Looks like you've been signed up!",
	}

	// send welcome email, increment billing period
	// send welcome email, increment billing period
	err = workflow.ExecuteActivity(ctx, activities.SendWelcomeEmail, data).Get(ctx, nil)
	if err != nil {
		return err
	}

	// start subscription period. execute until MaxBillingPeriods is reached
	for ; subscriptionPeriod < subscription.MaxSubscriptionPeriods; subscriptionPeriod++ {
		data := EmailInfo{
			EmailAddress: subscription.EmailInfo.EmailAddress,
			Mail:         "This is yet another email in the Subscription Workflow.",
		}

		err = workflow.ExecuteActivity(ctx, activities.SendSubscriptionEmail, data).Get(ctx, nil)
		if err != nil {
			return err
		}
		logger.Info("Sent content email", "EmailAddress", subscription.EmailInfo.EmailAddress)
		// Sleep the Workflow until the next subscription email needs to be sent.
		// This can be set to sleep every month between emails.
		if err = workflow.Sleep(ctx, duration); err != nil {
			return err
		}
	}
	return nil
}
// @@@SNIPEND
