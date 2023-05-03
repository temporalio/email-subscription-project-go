// @@@SNIPSTART subscription-workflow-go-workflow
package subscribe_emails

import (
	"errors"
	"strconv"
	"time"

	"go.temporal.io/sdk/workflow"
)

// Workflow definition
func SubscriptionWorkflow(ctx workflow.Context, subscription Subscription) error {
	// declare variables, duration, and logger
	var activities *Activities
	subscriptionPeriodNum := 0
	// duration can be set up to a month.
	duration := time.Minute

	logger := workflow.GetLogger(ctx)
	logger.Info("Subscription created for " + subscription.EmailInfo.EmailAddress)
	// Query result to be returned
	var queryResult string
	// Query handler
	e := workflow.SetQueryHandler(ctx, "GetDetails", func(input []byte) (string, error) {
		queryResult = subscription.EmailInfo.EmailAddress + " is on billing period " + strconv.Itoa(subscriptionPeriodNum) + " out of " + strconv.Itoa(subscription.Periods.MaxSubscriptionPeriods)
 		return queryResult, nil
	})
	if e != nil {
		logger.Info("SetQueryHandler failed: " + e.Error())
		return e 
	}

	var err error
	// set Activity Options. Timeout can be set to a longer timespan (such as a month)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
		WaitForCancellation: true,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	// Handle any cleanup, including cancellations.
	defer func() {
		if errors.Is(ctx.Err(), workflow.ErrCanceled) {
			newCtx, _ := workflow.NewDisconnectedContext(ctx)
			data := EmailInfo {
				EmailAddress: subscription.EmailInfo.EmailAddress,
				Mail:         "Oh my! Looks like your subscription has been canceled!",
			}
			// send cancellation email
			e := workflow.ExecuteActivity(newCtx, activities.SendCancellationEmail, data)
			if err != nil {
				logger.Error("Failed to send cancel email", "Error", e)
			} else {
				// Cancellation received, which will trigger an unsubscribe email.
				logger.Info("Sending cancellation email")
			}
			return
		}

		newCtx, _ := workflow.NewDisconnectedContext(ctx)
		// information for the cancelled subscription email
		data := EmailInfo {
				EmailAddress: subscription.EmailInfo.EmailAddress,
				Mail: "You have been unsubscribed from the Subscription Workflow. Good bye.",
		}
		logger.Info("Sending unsubscribe email to " + subscription.EmailInfo.EmailAddress)
		// send the cancelled subscription email
		err := workflow.ExecuteActivity(newCtx, activities.SendSubscriptionEndedEmail, data).Get(newCtx, nil)

		if err != nil {
			logger.Error("Unable to send unsubscribe message", "Error", err)
		}
	}()
	// handling for the first email ever
	logger.Info("Sending welcome email to " + subscription.EmailInfo.EmailAddress)

	data := EmailInfo {
				EmailAddress: subscription.EmailInfo.EmailAddress,
				Mail:         "Welcome! Looks like you've been signed up!",
		}
			
	// send welcome email, increment billing period
	err = workflow.ExecuteActivity(ctx, activities.SendWelcomeEmail, data).Get(ctx, nil)

	if err != nil {
		logger.Error("Failed to send welcome email", "Error", err)
	} else {
		subscriptionPeriodNum++
	}

	// start subscription period. execute until MaxBillingPeriods is reached
	for (subscriptionPeriodNum < subscription.Periods.MaxSubscriptionPeriods) {

		data := EmailInfo{
				EmailAddress: subscription.EmailInfo.EmailAddress,
				Mail:         "This is yet another email in the Subscription Workflow.",
		}

		err = workflow.ExecuteActivity(ctx, activities.SendSubscriptionEmail, data).Get(ctx, nil)

		if err != nil {
			logger.Error("Failed to send email ", "Error", err)
		}

		logger.Info("sent content email to " + subscription.EmailInfo.EmailAddress)

		// increment billing period for successful email
		subscriptionPeriodNum++
		// Sleep the Workflow until the next subscription email needs to be sent.
		// This can be set to sleep every month between emails.
		workflow.Sleep(ctx, duration)
	}

	return err
}
// @@@SNIPEND

