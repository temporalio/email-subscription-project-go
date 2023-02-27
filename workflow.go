package subscribe_emails

import (
	"errors"
	"time"

	"go.temporal.io/sdk/workflow"
)

// Workflow definition
func SubscriptionWorkflow(ctx workflow.Context, subscription Subscription) {
	// set up logger
	logger := workflow.GetLogger(ctx)
	logger.Info("Subscription created for " + subscription.EmailAddress)

	// set up Activity Options
	duration := time.Minute
	ao := workflow.ActivityOptions {
		StartToCloseTimeout: 10 * time.Minute,
		WaitForCancellation: true,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	// handle subscription
	defer func() {
		if !errors.Is(ctx.Err(), workflow.ErrCanceled) {
			return
		}

		// Cancellation

		newCtx, _ := workflow.NewDisconnectedContext(ctx)
		
		data := EmailInfo {
			EmailAddress: subscription.EmailAddress,
			Mail: subscription.Campaign.UnsubscribeEmail,
		}

		logger.Info("Sending unsubscribe email to " + subscription.EmailAddress)

		err := workflow.ExecuteActivity(newCtx, SendContentEmail, data).Get(newCtx, nil)

		if err != nil {
			logger.Error("Unable to send unsubscribe message", "Error", err)
		}
	}()

	logger.Info("Sending welcome email to " + subscription.EmailAddress)

	data := EmailInfo {
		EmailAddress: subscription.EmailAddress,
		Mail: subscription.Campaign.WelcomeEmail,
	}

	err := workflow.ExecuteActivity(ctx, SendContentEmail, data).Get(ctx, nil)

	if err != nil {
		logger.Error("Failed to send welcome email", "Error", err)
	}

	for _, mail := range subscription.Campaign.Mails {
		data := EmailInfo {
			EmailAddress: subscription.EmailAddress,
			Mail: 		  mail,
		}

		err = workflow.ExecuteActivity(ctx, SendContentEmail, data).Get(ctx, nil)

		if err != nil {
			logger.Error("Failed to send email " +mail, "Error", err)
		}

		logger.Info("Sent content email " + mail + " to " + subscription.EmailAddress)

		workflow.Sleep(ctx, duration)
	}

}
