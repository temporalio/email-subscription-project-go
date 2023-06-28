// @@@SNIPSTART subscription-workflow-go-main
package subscribeemails

import (
	"errors"
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

// Workflow definition
func SubscriptionWorkflow(ctx workflow.Context, emailDetails EmailDetails) error {
	duration := 12 * time.Second
	logger := workflow.GetLogger(ctx)
	logger.Info("Subscription created", "EmailAddress", emailDetails.EmailAddress)
	// Query handler
	err := workflow.SetQueryHandler(ctx, "GetDetails", func() (string, error) {
		return fmt.Sprintf("%v is on email #%v ",
			emailDetails.EmailAddress,
			emailDetails.SubscriptionCount), nil
	})
	if err != nil {
		return err
	}
	// variable for Activity Options. Timeout can be set to a longer timespan (such as a month)
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
		WaitForCancellation: true,
	})

	// Handle any cleanup, including cancellations.
	defer func() {
		newCtx, cancel := workflow.NewDisconnectedContext(ctx)
		defer cancel()

		if errors.Is(ctx.Err(), workflow.ErrCanceled) {
			data := EmailDetails{
				EmailAddress:      emailDetails.EmailAddress,
				Message:           "Your subscription has been canceled. Sorry to see you go!",
				IsSubscribed:      false,
				SubscriptionCount: emailDetails.SubscriptionCount,
			}
			// send cancellation email
			err := workflow.ExecuteActivity(newCtx, SendEmail, data).Get(newCtx, nil)
			if err != nil {
				logger.Error("Failed to send cancellation email", "Error", err)
			} else {
				// Cancellation received.
				logger.Info("Sent cancellation email", "EmailAddress", emailDetails.EmailAddress)
			}
		}
	}()

	// handling for the first email ever
	logger.Info("Sending welcome email", "EmailAddress", emailDetails.EmailAddress)
	emailDetails.SubscriptionCount++
	data := EmailDetails{
		EmailAddress:      emailDetails.EmailAddress,
		Message:           "Welcome! Looks like you've been signed up!",
		IsSubscribed:      true,
		SubscriptionCount: emailDetails.SubscriptionCount,
	}

	// send welcome email, increment billing period
	err = workflow.ExecuteActivity(ctx, SendEmail, data).Get(ctx, nil)
	if err != nil {
		return err
	}

	// start subscription period. execute until no longer subscribed
	for emailDetails.IsSubscribed {
		emailDetails.SubscriptionCount++
		data := EmailDetails{
			EmailAddress:      emailDetails.EmailAddress,
			Message:           "This is yet another email in the Subscription Workflow.",
			IsSubscribed:      true,
			SubscriptionCount: emailDetails.SubscriptionCount,
		}

		err = workflow.ExecuteActivity(ctx, SendEmail, data).Get(ctx, nil)
		if err != nil {
			return err
		}
		logger.Info("Sent content email", "EmailAddress", emailDetails.EmailAddress)
		// Sleep the Workflow until the next subscription email needs to be sent.
		// This can be set to sleep every month between emails.
		if err = workflow.Sleep(ctx, duration); err != nil {
			return err
		}
	}
	return nil
}

// @@@SNIPEND
