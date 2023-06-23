// @@@SNIPSTART subscription-workflow-go-workflow
package subscribeemails

import (
	"errors"
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// Workflow definition
func SubscriptionWorkflow(ctx workflow.Context, emailDetails EmailDetails) error {
	//subscriptionPeriodCount := emailDetails.SubscriptionCount
	duration := 5 * time.Second
	logger := workflow.GetLogger(ctx)
	logger.Info("Subscription created", "EmailAddress", emailDetails.EmailAddress)
	// Query handler
	err := workflow.SetQueryHandler(ctx, "GetDetails", func(input []byte) (string, error) {
		return fmt.Sprintf("%v is on subscription period %v out of %v",
			emailDetails.EmailAddress,
			emailDetails.SubscriptionCount,
			emailDetails.MaxSubscriptionPeriods), nil
	})
	if err != nil {
		return err
	}
	// variable for Activity Options. Timeout can be set to a longer timespan (such as a month)
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
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

		if errors.Is(ctx.Err(), workflow.ErrCanceled) && !emailDetails.IsSubscribed {
			data := EmailDetails{
				EmailAddress:           emailDetails.EmailAddress,
				Message:                "Your subscription has been canceled. Sorry to see you go!",
				IsSubscribed:           false,
				SubscriptionCount:      emailDetails.SubscriptionCount,
				MaxSubscriptionPeriods: emailDetails.MaxSubscriptionPeriods,
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

		// information for the newly-ended subscription email
		if emailDetails.MaxSubscriptionPeriods == emailDetails.SubscriptionCount {
			data := EmailDetails{
				EmailAddress:           emailDetails.EmailAddress,
				Message:                "You have been unsubscribed from the Subscription Workflow. Goodbye.",
				IsSubscribed:           false,
				MaxSubscriptionPeriods: emailDetails.MaxSubscriptionPeriods,
				SubscriptionCount:      emailDetails.SubscriptionCount,
			}
			logger.Info("Sending unsubscribe email", "EmailAddress", emailDetails.EmailAddress)
			// send the cancelled subscription email
			err := workflow.ExecuteActivity(newCtx, SendEmail, data).Get(newCtx, nil)

			if err != nil {
				logger.Error("Unable to send unsubscribe message", "Error", err)
			}
		}
	}()

	// handling for the first email ever
	logger.Info("Sending welcome email", "EmailAddress", emailDetails.EmailAddress)
	emailDetails.SubscriptionCount++
	data := EmailDetails{
		EmailAddress:           emailDetails.EmailAddress,
		Message:                "Welcome! Looks like you've been signed up!",
		IsSubscribed:           true,
		MaxSubscriptionPeriods: emailDetails.MaxSubscriptionPeriods,
		SubscriptionCount:      emailDetails.SubscriptionCount,
	}

	// send welcome email, increment billing period
	err = workflow.ExecuteActivity(ctx, SendEmail, data).Get(ctx, nil)
	if err != nil {
		return err
	}

	// start subscription period. execute until MaxBillingPeriods is reached
	for emailDetails.SubscriptionCount < emailDetails.MaxSubscriptionPeriods {
		emailDetails.SubscriptionCount++
		data := EmailDetails{
			EmailAddress:           emailDetails.EmailAddress,
			Message:                "This is yet another email in the Subscription Workflow.",
			IsSubscribed:           true,
			MaxSubscriptionPeriods: emailDetails.MaxSubscriptionPeriods,
			SubscriptionCount:      emailDetails.SubscriptionCount,
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
